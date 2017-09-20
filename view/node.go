package view

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/full"
	"gomatcha.io/matcha/paint"
	pb "gomatcha.io/matcha/pb/view"
)

var maxId int64

// Middleware is called on the result of View.Build(*context).
type middleware interface {
	Build(Context, *Model)
	MarshalProtobuf() proto.Message
	Key() string
}

// root contains your view hierarchy.
type root struct {
	id     int64
	root   *nodeRoot
	size   layout.Point
	ticker *internal.Ticker
}

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/view NewRoot", func(v View) *root {
		return _newRoot(v)
	})
}

// NewRoot initializes a Root with screen s.
func _newRoot(v View) *root {
	r := &root{
		root: newRoot(v),
		id:   atomic.AddInt64(&maxId, 1),
	}
	r.start()
	return r
}

func (r *root) start() {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	if r.ticker != nil {
		return
	}

	id := r.id
	r.ticker = internal.NewTicker(time.Hour * 99999)
	_ = r.ticker.Notify(func() {
		matcha.MainLocker.Lock()
		defer matcha.MainLocker.Unlock()

		if !r.root.update(r.size) {
			// nothing changed
			return
		}

		pb, err := r.root.MarshalProtobuf2()
		if err != nil {
			fmt.Println("err", err)
			return
		}

		// fmt.Println(r.root.node.debugString())
		fmt.Println("Update") // TODO(KD): Remove.
		if runtime.GOOS == "android" {
			bridge.Bridge("").Call("updateViewWithProtobuf", bridge.Int64(id), bridge.Bytes(pb))
		} else if runtime.GOOS == "darwin" {
			bridge.Bridge("").Call("updateId:withProtobuf:", bridge.Int64(id), bridge.Bytes(pb))
		}
	})
}

func (r *root) Stop() {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	if r.ticker == nil {
		return
	}
	r.ticker.Stop()
}

func (r *root) Call(funcId string, viewId int64, args []reflect.Value) []reflect.Value {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.root.call(funcId, viewId, args)
}

// Id returns the unique identifier for r.
func (r *root) Id() int64 {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.id
}

func (r *root) ViewId() Id {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.root.node.id
}

// Size returns the size of r.
func (r *root) Size() layout.Point {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.size
}

// SetSize sets the size of r.
func (r *root) SetSize(width, height float64) {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	r.size = layout.Pt(width, height)
	r.root.addFlag(r.root.node.id, layoutFlag)
}

func newId() Id {
	return Id(atomic.AddInt64(&maxId, 1))
}

type Context interface {
	Path() []Id
}

// viewContext specifies the supporting context for building a View.
type viewContext struct {
	valid bool
	node  *node
}

// Path returns the path of Ids from the root to the view.
func (ctx *viewContext) Path() []Id {
	if ctx.node == nil {
		return []Id{0}
	}
	return ctx.node.path
}

type updateFlag int

const (
	buildFlag updateFlag = 1 << iota
	layoutFlag
	paintFlag
)

func (f updateFlag) needsBuild() bool {
	return f&buildFlag != 0
}

func (f updateFlag) needsLayout() bool {
	return f&buildFlag != 0 || f&layoutFlag != 0
}

func (f updateFlag) needsPaint() bool {
	return f&buildFlag != 0 || f&layoutFlag != 0 || f&paintFlag != 0
}

type nodeRoot struct {
	node        *node
	nodes       map[Id]*node
	middlewares []middleware

	flagMu      sync.Mutex
	updateFlags map[Id]updateFlag
}

func newRoot(v View) *nodeRoot {
	id := newId()
	root := &nodeRoot{}
	root.node = &node{
		id:   id,
		path: []Id{id},
		view: v,
		root: root,
	}
	root.updateFlags = map[Id]updateFlag{id: buildFlag}
	for _, i := range internal.Middlewares() {
		root.middlewares = append(root.middlewares, i().(middleware))
	}
	return root
}

func (root *nodeRoot) addFlag(id Id, f updateFlag) {
	root.flagMu.Lock()
	defer root.flagMu.Unlock()

	root.updateFlags[id] |= f
}

func (root *nodeRoot) update(size layout.Point) bool {
	root.flagMu.Lock()
	defer root.flagMu.Unlock()

	var flag updateFlag
	for _, v := range root.updateFlags {
		flag |= v
	}

	updated := false
	if flag.needsBuild() {
		root.build()
		updated = true
	}
	if flag.needsLayout() {
		root.layout(size, size)
		updated = true
	}
	if flag.needsPaint() {
		root.paint()
		updated = true
	}
	root.updateFlags = map[Id]updateFlag{}
	return updated
}

func (root *nodeRoot) MarshalProtobuf2() ([]byte, error) {
	return proto.Marshal(root.MarshalProtobuf())
}

func (root *nodeRoot) MarshalProtobuf() *pb.Root {
	m := map[int64]*pb.LayoutPaintNode{}
	root.node.marshalLayoutPaintProtobuf(m)

	m2 := map[int64]*pb.BuildNode{}
	root.node.marshalBuildProtobuf(m2)

	m3 := map[string]*any.Any{}
	for _, i := range root.middlewares {
		var _any *any.Any
		if a, err := ptypes.MarshalAny(i.MarshalProtobuf()); err == nil {
			_any = a
		}
		m3[i.Key()] = _any
	}

	return &pb.Root{
		LayoutPaintNodes: m,
		BuildNodes:       m2,
		Middleware:       m3,
	}
}

func (root *nodeRoot) build() {
	root.nodes = map[Id]*node{
		root.node.id: root.node,
	}

	// Rebuild
	root.node.build()
}

func (root *nodeRoot) layout(minSize layout.Point, maxSize layout.Point) {
	g := root.node.layout(minSize, maxSize)
	g.Frame = g.Frame.Add(layout.Pt(-g.Frame.Min.X, -g.Frame.Min.Y)) // Move Frame.Min to the origin.
	root.node.layoutGuide = &g
}

func (root *nodeRoot) paint() {
	root.node.paint()
}

func (root *nodeRoot) call(funcId string, viewId int64, args []reflect.Value) []reflect.Value {
	node, ok := root.nodes[Id(viewId)]
	if !ok || node.model == nil {
		fmt.Println("root.call(): no node found", ok, node.model)
		return nil
	}

	f, ok := node.model.NativeFuncs[funcId]
	if !ok {
		fmt.Println("root.call(): no func found", funcId, node.model.NativeFuncs)
		return nil
	}
	v := reflect.ValueOf(f)

	return v.Call(args)
}

type node struct {
	id    Id
	path  []Id
	root  *nodeRoot
	view  View
	stage Stage

	buildId       int64
	buildPbId     int64
	buildNotify   bool
	buildNotifyId comm.Id
	model         *Model
	children      []*node

	layoutId       int64
	layoutNotify   bool
	layoutNotifyId comm.Id
	layoutGuide    *layout.Guide
	layoutMinSize  layout.Point
	layoutMaxSize  layout.Point

	paintId       int64
	paintNotify   bool
	paintNotifyId comm.Id
	paintOptions  paint.Style
}

func (n *node) marshalLayoutPaintProtobuf(m map[int64]*pb.LayoutPaintNode) {
	guide := n.layoutGuide
	if n.layoutGuide == nil {
		guide = &layout.Guide{}
		fmt.Println("View is missing layout guide", n.id, n.view)
	}

	// Sort children by zIndex for performance reasons.
	childOrder := []struct {
		id int64
		z  int
	}{}
	for _, i := range n.children {
		z := 0
		if i.layoutGuide != nil {
			z = i.layoutGuide.ZIndex
		}
		childOrder = append(childOrder, struct {
			id int64
			z  int
		}{id: int64(i.id), z: z})
	}
	sort.Slice(childOrder, func(i, j int) bool {
		return childOrder[i].z < childOrder[j].z
	})
	order := []int64{}
	for _, i := range childOrder {
		order = append(order, i.id)
	}

	m[int64(n.id)] = &pb.LayoutPaintNode{
		Id:       int64(n.id),
		LayoutId: n.layoutId,
		PaintId:  n.paintId,

		// LayoutGuide: guide.MarshalProtobuf(),
		Minx:       guide.Frame.Min.X,
		Miny:       guide.Frame.Min.Y,
		Maxx:       guide.Frame.Max.X,
		Maxy:       guide.Frame.Max.Y,
		ZIndex:     int64(guide.ZIndex),
		ChildOrder: order,

		PaintStyle: n.paintOptions.MarshalProtobuf(),
	}
	for _, v := range n.children {
		v.marshalLayoutPaintProtobuf(m)
	}
}

func (n *node) marshalBuildProtobuf(m map[int64]*pb.BuildNode) {
	for _, v := range n.children {
		v.marshalBuildProtobuf(m)
	}

	// Don't build if nothing has changed
	if n.buildPbId == n.buildId {
		return
	}
	n.buildPbId = n.buildId

	children := []int64{}
	for _, v := range n.children {
		children = append(children, int64(v.id))
	}

	var nativeViewState *any.Any
	if a, err := ptypes.MarshalAny(n.model.NativeViewState); err == nil {
		nativeViewState = a
	}

	nativeValues := map[string]*google_protobuf.Any{}
	for k, v := range n.model.NativeOptions {
		a, err := ptypes.MarshalAny(v)
		if err != nil {
			fmt.Println("Error enocding native value: ", err)
			continue
		}
		nativeValues[k] = a
	}

	m[int64(n.id)] = &pb.BuildNode{
		Id:          int64(n.id),
		BuildId:     n.buildId,
		Children:    children,
		BridgeName:  n.model.NativeViewName,
		BridgeValue: nativeViewState,
		Values:      nativeValues,
	}
}

func (n *node) build() {
	if n.root.updateFlags[n.id].needsBuild() {
		n.buildId += 1

		// Send lifecycle event to new children.
		if n.stage == StageDead {
			n.view.Lifecycle(n.stage, StageVisible)
			n.stage = StageVisible
		}

		// Generate the new viewModel.
		ctx := &viewContext{valid: true, node: n}
		temp := n.view.Build(ctx)
		viewModel := &temp

		// Call middleware
		for _, i := range n.root.middlewares {
			i.Build(ctx, viewModel)
		}
		ctx.valid = false

		//
		prevChildren := make([]*node, len(n.children))
		copy(prevChildren, n.children)

		children := []*node{}
		for _, i := range viewModel.Children {
			// Find the corresponding previous node.
			var prevNode *node

			iKey := i.ViewKey()
			iName := internal.ReflectName(i)
			for jIdx, j := range prevChildren {
				jName := internal.ReflectName(j.view)
				jKey := j.view.ViewKey()

				if iKey == jKey && iName == jName {
					prevNode = j

					// delete from prevchildren
					copy(prevChildren[jIdx:], prevChildren[jIdx+1:])
					prevChildren[len(prevChildren)-1] = nil
					prevChildren = prevChildren[:len(prevChildren)-1]
					break
				}
			}

			if prevNode != nil {
				// If view was modified...
				prevView := prevNode.view
				newView := i

				// Copy all public fields from new to old, if embed.Update is called.
				if prevView != newView {
					embedUpdate = false
					prevView.Update(newView)
					if embedUpdate {
						CopyFields(prevView, newView)
					}
				}

				// Add in the previous node.
				children = append(children, prevNode)

				// Mark as needing rebuild
				n.root.updateFlags[prevNode.id] |= buildFlag
			} else {
				// If view was added for the first time...
				newView := i
				id := newId()

				// Add in a new node.
				path := make([]Id, len(n.path)+1)
				copy(path, n.path)
				path[len(n.path)] = id
				children = append(children, &node{
					id:   id,
					path: path,
					view: newView,
					root: n.root,
				})

				// Mark as needing rebuild
				n.root.updateFlags[id] |= buildFlag
			}
		}

		// Send lifecycle event to removed childern.
		for _, i := range prevChildren {
			i.done()
		}

		// Watch for build changes, if we haven't
		if !n.buildNotify {
			n.buildNotifyId = n.view.Notify(func() {
				n.root.addFlag(n.id, buildFlag)
			})
			n.buildNotify = true
		}

		// Watch for layout changes.
		if n.layoutNotify {
			n.model.Layouter.Unnotify(n.layoutNotifyId)
			n.layoutNotify = false
		}
		if viewModel.Layouter != nil {
			n.layoutNotifyId = viewModel.Layouter.Notify(func() {
				n.root.addFlag(n.id, layoutFlag)
			})
			n.layoutNotify = true
		}

		// Watch for paint changes.
		if n.paintNotify {
			n.model.Painter.Unnotify(n.paintNotifyId)
			n.paintNotify = false
		}
		if viewModel.Painter != nil {
			n.paintNotifyId = viewModel.Painter.Notify(func() {
				n.root.addFlag(n.id, paintFlag)
			})
			n.paintNotify = true
		}

		n.children = children

		n.model = viewModel
	}

	// Recursively update children.
	for _, i := range n.children {
		i.build()

		// Also add to the root
		n.root.nodes[i.id] = i
	}
}

func (n *node) layout(minSize layout.Point, maxSize layout.Point) layout.Guide {
	n.layoutId += 1

	// If node has no children, has the same min/max size, and does not need relayout, return the previous guide.
	if len(n.children) == 0 && n.layoutGuide != nil && n.layoutMinSize == minSize && n.layoutMaxSize == maxSize && !n.root.updateFlags[n.id].needsLayout() {
		return *n.layoutGuide
	}
	n.layoutMinSize = minSize
	n.layoutMaxSize = maxSize

	// Create the LayoutContext
	ctx := &layoutContext{
		minSize:    minSize,
		maxnSize:   maxSize,
		childCount: len(n.children),
		layoutFunc: func(idx int, minSize, maxSize layout.Point) layout.Guide {
			if idx >= len(n.children) {
				fmt.Println("Attempting to layout unknown child: ", idx)
				return layout.Guide{}
			}
			child := n.children[idx]
			return child.layout(minSize, maxSize)
		},
	}

	// Perform layout
	layouter := n.model.Layouter
	if layouter == nil {
		layouter = &full.Layouter{}
	}
	g, gs := layouter.Layout(ctx)
	g = ctx.fitGuide(g)

	//
	for idx, i := range n.children {
		g2 := gs[idx]
		i.layoutGuide = &g2
	}
	return g
}

func (n *node) paint() {
	if n.root.updateFlags[n.id].needsPaint() {
		n.paintId += 1

		if p := n.model.Painter; p != nil {
			n.paintOptions = p.PaintStyle()
		} else {
			n.paintOptions = paint.Style{}
		}
	}

	// Recursively update children
	for _, v := range n.children {
		v.paint()
	}
}

func (n *node) done() {
	n.view.Lifecycle(n.stage, StageDead)
	n.stage = StageDead

	if n.buildNotify {
		n.view.Unnotify(n.buildNotifyId)
	}
	if n.layoutNotify {
		n.model.Layouter.Unnotify(n.layoutNotifyId)
	}
	if n.paintNotify {
		n.model.Painter.Unnotify(n.paintNotifyId)
	}

	for _, i := range n.children {
		i.done()
	}
}

func (n *node) debugString() string {
	all := []string{}
	for _, i := range n.children {
		lines := strings.Split(i.debugString(), "\n")
		for idx, line := range lines {
			lines[idx] = "|	" + line
		}
		all = append(all, lines...)
	}

	str := fmt.Sprintf("{%p Id:%v,%v,%v,%v View:%v Node:%p Layout:%v}", n, n.id, n.buildId, n.layoutId, n.paintId, n.view, n.model, n.layoutGuide.Frame)
	if len(all) > 0 {
		str += "\n" + strings.Join(all, "\n")
	}
	return str
}

type layoutContext struct {
	minSize    layout.Point
	maxnSize   layout.Point
	childCount int
	layoutFunc func(int, layout.Point, layout.Point) layout.Guide // TODO(KD): this should be private...
}

func (l *layoutContext) MinSize() layout.Point {
	return l.minSize
}

func (l *layoutContext) MaxSize() layout.Point {
	return l.maxnSize
}

func (l *layoutContext) ChildCount() int {
	return l.childCount
}

func (l *layoutContext) LayoutChild(idx int, minSize, maxSize layout.Point) layout.Guide {
	g := l.layoutFunc(idx, minSize, maxSize)
	g.Frame = g.Frame.Add(layout.Pt(-g.Frame.Min.X, -g.Frame.Min.Y))
	return g
}

func (l *layoutContext) fitGuide(g layout.Guide) layout.Guide {
	if g.Width() < l.MinSize().X {
		g.Frame.Max.X = l.MinSize().X - g.Frame.Min.X
	}
	if g.Height() < l.MinSize().Y {
		g.Frame.Max.Y = l.MinSize().Y - g.Frame.Min.Y
	}
	if g.Width() > l.MaxSize().X {
		g.Frame.Max.X = l.MaxSize().X - g.Frame.Min.X
	}
	if g.Height() > l.MaxSize().Y {
		g.Frame.Max.Y = l.MaxSize().Y - g.Frame.Min.Y
	}
	return g
}

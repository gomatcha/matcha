package view

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/full"
	"gomatcha.io/matcha/paint"
	pb "gomatcha.io/matcha/proto/view"
)

var maxId int64

// Middleware is called on the result of View.Build(Context).
type middleware interface {
	Build(Context, *Model)
	MarshalProtobuf() proto.Message
	Key() string
}

// root contains your view hierarchy.
type root struct {
	id         int64
	root       *nodeRoot
	size       layout.Point
	ticker     *internal.Ticker
	printDebug bool
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

	id := r.id
	r.ticker = internal.NewTicker(func() {
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

		if r.printDebug {
			fmt.Println(r.root.recursiveString())
		}

		success := false
		if runtime.GOOS == "android" {
			success = bridge.Bridge("").Call("updateViewWithProtobuf", bridge.Int64(id), bridge.Bytes(pb)).ToBool()
		} else if runtime.GOOS == "darwin" {
			success = bridge.Bridge("").Call("updateId:withProtobuf:", bridge.Int64(id), bridge.Bytes(pb)).ToBool()
		}
		if !success {
			r.ticker.Stop()
		}
	})
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

func (r *root) PrintDebug() {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	fmt.Println(r.root.recursiveString())
}

func (r *root) SetPrintDebug(v bool) {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	r.printDebug = v
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
	printDebug  bool
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

	root.printDebug = false

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

	if root.printDebug {
		fmt.Println(root.recursiveString())
	}

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
	root.node.layoutGuide = g
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

func (root *nodeRoot) recursiveString() string {
	lines := strings.Split(root.node.recursiveString(), "\n")
	for idx, line := range lines {
		lines[idx] = "    " + line
	}
	all := append([]string{"View hierarchy:"}, lines...)
	return strings.Join(all, "\n")
}

type node struct {
	id    Id
	path  []Id
	root  *nodeRoot
	view  View
	stage Stage

	buildId         int64
	buildIsNotified bool
	buildNotifyId   comm.Id
	model           *Model
	children        []*node

	layoutId          int64
	layoutIsNotified  bool
	layoutNotifyId    comm.Id
	layoutGuide       layout.Guide
	layoutMinSize     layout.Point
	layoutMaxSize     layout.Point
	layoutDebugString string

	paintId         int64
	paintIsNotified bool
	paintNotifyId   comm.Id
	paintOptions    paint.Style

	// Skip unneccesary serialization
	lastBuildId  int64
	lastLayoutId int64
	lastPaintId  int64
}

func (n *node) marshalLayoutPaintProtobuf(m map[int64]*pb.LayoutPaintNode) {
	// Marshal children
	for _, v := range n.children {
		v.marshalLayoutPaintProtobuf(m)
	}

	// Don't serialize if nothing has changed
	if n.lastLayoutId == n.layoutId && n.lastPaintId == n.paintId {
		return
	}
	n.lastLayoutId = n.layoutId
	n.lastPaintId = n.paintId

	// Sort children by zIndex for performance reasons.
	childOrder := []struct {
		id int64
		z  int
	}{}
	for _, i := range n.children {
		childOrder = append(childOrder, struct {
			id int64
			z  int
		}{id: int64(i.id), z: i.layoutGuide.ZIndex})
	}
	sort.Slice(childOrder, func(i, j int) bool {
		return childOrder[i].z < childOrder[j].z
	})
	order := []int64{}
	for _, i := range childOrder {
		order = append(order, i.id)
	}

	guide := n.layoutGuide
	s := n.paintOptions
	lpnode := &pb.LayoutPaintNode{
		Id:       int64(n.id),
		LayoutId: n.layoutId,
		PaintId:  n.paintId,

		Minx:       guide.Frame.Min.X,
		Miny:       guide.Frame.Min.Y,
		Maxx:       guide.Frame.Max.X,
		Maxy:       guide.Frame.Max.Y,
		ZIndex:     int64(guide.ZIndex),
		ChildOrder: order,

		Transparency:       s.Transparency,
		HasBackgroundColor: s.BackgroundColor != nil,
		HasBorderColor:     s.BorderColor != nil,
		BorderWidth:        s.BorderWidth,
		CornerRadius:       s.CornerRadius,
		ShadowRadius:       s.ShadowRadius,
		ShadowOffsetX:      s.ShadowOffset.X,
		ShadowOffsetY:      s.ShadowOffset.Y,
		HasShadowColor:     s.ShadowColor != nil,
	}
	if s.BackgroundColor != nil {
		r, g, b, a := s.BackgroundColor.RGBA()
		lpnode.HasBackgroundColor = true
		lpnode.BackgroundColorRed = r
		lpnode.BackgroundColorGreen = g
		lpnode.BackgroundColorBlue = b
		lpnode.BackgroundColorAlpha = a
	}
	if s.BorderColor != nil {
		r, g, b, a := s.BorderColor.RGBA()
		lpnode.HasBorderColor = true
		lpnode.BorderColorRed = r
		lpnode.BorderColorGreen = g
		lpnode.BorderColorBlue = b
		lpnode.BorderColorAlpha = a
	}
	if s.ShadowColor != nil {
		r, g, b, a := s.ShadowColor.RGBA()
		lpnode.HasShadowColor = true
		lpnode.ShadowColorRed = r
		lpnode.ShadowColorGreen = g
		lpnode.ShadowColorBlue = b
		lpnode.ShadowColorAlpha = a
	}
	m[int64(n.id)] = lpnode
}

func (n *node) marshalBuildProtobuf(m map[int64]*pb.BuildNode) {
	// Marshal children
	for _, v := range n.children {
		v.marshalBuildProtobuf(m)
	}

	// Don't serialize if nothing has changed
	if n.lastBuildId == n.buildId {
		return
	}
	n.lastBuildId = n.buildId

	children := []int64{}
	for _, v := range n.children {
		children = append(children, int64(v.id))
	}

	nativeValues := map[string][]byte{}
	for k, v := range n.model.NativeOptions {
		nativeValues[k] = v
	}

	m[int64(n.id)] = &pb.BuildNode{
		Id:          int64(n.id),
		BuildId:     n.buildId,
		Children:    children,
		BridgeName:  n.model.NativeViewName,
		BridgeValue: n.model.NativeViewState,
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
		if !n.buildIsNotified {
			n.buildNotifyId = n.view.Notify(func() {
				n.root.addFlag(n.id, buildFlag)
			})
			n.buildIsNotified = true
		}

		// Watch for layout changes.
		if n.layoutIsNotified {
			n.model.Layouter.Unnotify(n.layoutNotifyId)
			n.layoutIsNotified = false
		}
		if viewModel.Layouter != nil {
			n.layoutNotifyId = viewModel.Layouter.Notify(func() {
				n.root.addFlag(n.id, layoutFlag)
			})
			n.layoutIsNotified = true
		}

		// Watch for paint changes.
		if n.paintIsNotified {
			n.model.Painter.Unnotify(n.paintNotifyId)
			n.paintIsNotified = false
		}
		if viewModel.Painter != nil {
			n.paintNotifyId = viewModel.Painter.Notify(func() {
				n.root.addFlag(n.id, paintFlag)
			})
			n.paintIsNotified = true
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
	// If node has no children, has the same min/max size, and does not need relayout, return the previous guide.
	if len(n.children) == 0 && n.layoutMinSize == minSize && n.layoutMaxSize == maxSize && !n.root.updateFlags[n.id].needsLayout() {
		return n.layoutGuide
	}
	// If node has no children, and min/max size are equivalent, return the min size.
	if len(n.children) == 0 && minSize == maxSize {
		return layout.Guide{Frame: layout.Rt(0, 0, minSize.X, minSize.Y)}
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

	// Assign guides to children
	for idx, i := range n.children {
		childGuide := gs[idx]
		if !isRectValid(childGuide.Frame) {
			fmt.Printf("Invalid rect for view. Rect:%v View:%v\n", childGuide.Frame, i)
			n.root.printDebug = true
		}
		if i.layoutGuide != childGuide {
			i.layoutGuide = childGuide
			i.layoutId += 1
		}
	}
	return g
}

func (n *node) paint() {
	if n.root.updateFlags[n.id].needsPaint() {
		opts := paint.Style{}
		if p := n.model.Painter; p != nil {
			opts = p.PaintStyle()
		}

		if opts != n.paintOptions {
			n.paintId += 1
			n.paintOptions = opts
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

	if n.buildIsNotified {
		n.view.Unnotify(n.buildNotifyId)
	}
	if n.layoutIsNotified {
		n.model.Layouter.Unnotify(n.layoutNotifyId)
	}
	if n.paintIsNotified {
		n.model.Painter.Unnotify(n.paintNotifyId)
	}

	ctx := &viewContext{valid: true, node: n}
	for _, i := range n.root.middlewares {
		i.Build(ctx, nil)
	}
	ctx.valid = false

	for _, i := range n.children {
		i.done()
	}
}

func (n *node) String() string {
	viewLine := "{" + reflect.TypeOf(n.view).String() + " Id:" + strconv.Itoa(int(n.id)) + " "
	split := strings.SplitN(fmt.Sprintf("%+v", n.view), "lastField:{}} ", 2)
	if len(split) == 2 {
		viewLine += split[1]
	} else {
		viewLine += strings.TrimPrefix("{", split[0])
	}
	return viewLine
}

func (n *node) recursiveString() string {
	// View line
	viewLine := reflect.TypeOf(n.view).String() + " - {Id:" + strconv.Itoa(int(n.id)) + " "
	split := strings.SplitN(fmt.Sprintf("%+v", n.view), "lastField:{}} ", 2)
	if len(split) == 2 {
		viewLine += split[1]
	} else {
		viewLine += strings.TrimPrefix("{", split[0])
	}

	// Layout line
	width := ""
	if n.layoutMinSize.X == n.layoutMaxSize.X {
		width = fmt.Sprintf("W=%v", n.layoutMinSize.X)
	} else {
		width = fmt.Sprintf("%v<W<%v", n.layoutMinSize.X, n.layoutMaxSize.X)
	}

	height := ""
	if n.layoutMinSize.Y == n.layoutMaxSize.Y {
		height = fmt.Sprintf("H=%v", n.layoutMinSize.Y)
	} else {
		height = fmt.Sprintf("%v<H<%v", n.layoutMinSize.Y, n.layoutMaxSize.Y)
	}

	layout := ""
	var childLayouts []string
	if ld, ok := n.model.Layouter.(layouterDebug); ok {
		layout, childLayouts = ld.DebugStrings()
	}

	for idx, i := range n.children {
		childLayout := ""
		if childLayouts != nil {
			childLayout = fmt.Sprintf("{%T %v}", n.model.Layouter, childLayouts[idx])
		} else {
			childLayout = fmt.Sprintf("{%T}", n.model.Layouter)
		}
		i.layoutDebugString = childLayout
	}

	nodeLayout := ""
	if layout != "" {
		nodeLayout = fmt.Sprintf("{%T %v %v %v}", n.model.Layouter, width, height, layout)
	} else {
		nodeLayout = fmt.Sprintf("{%T %v %v}", n.model.Layouter, width, height)
	}
	parentLayout := n.layoutDebugString
	if parentLayout == "" {
		parentLayout = "{<nil>}"
	}
	layoutLine := "\n|Layout:" + parentLayout + "->" + nodeLayout + "->" + n.layoutGuide.Frame.String()

	// Paint line
	paintLine := ""
	if n.model.Painter != nil {
		paintLine = "\n|Paint:" + n.paintOptions.String()
	}

	// Options line
	optionsLine := ""
	if len(n.model.Options) != 0 {
		optionsLine = "\n|Options:" + fmt.Sprintf("{%+v}", n.model.Options)
	}

	// Build string
	str := viewLine + layoutLine + paintLine + optionsLine

	// Build children
	all := []string{}
	for _, i := range n.children {
		lines := strings.Split(i.recursiveString(), "\n")
		for idx, line := range lines {
			lines[idx] = "|    " + line
		}
		all = append(all, lines...)
	}
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

type layouterDebug interface {
	DebugStrings() (string, []string)
}

func isRectValid(r layout.Rect) bool {
	if math.IsInf(r.Min.X, 0) || math.IsNaN(r.Min.X) ||
		math.IsInf(r.Min.Y, 0) || math.IsNaN(r.Min.Y) ||
		math.IsInf(r.Max.X, 0) || math.IsNaN(r.Max.X) ||
		math.IsInf(r.Max.Y, 0) || math.IsNaN(r.Max.Y) {
		return false
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return false
	}
	return true
}

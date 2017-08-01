package view

import (
	"fmt"
	"reflect"
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
	Build(*Context, *Model)
	MarshalProtobuf() proto.Message
	Key() string
}

type Root struct {
	id     int64
	root   *root
	size   layout.Point
	ticker *internal.Ticker
}

// NewRoot initializes a Root with screen s.
func NewRoot(v View) *Root {
	r := &Root{
		root: newRoot(v),
		id:   atomic.AddInt64(&maxId, 1),
	}
	r.start()
	return r
}

func (r *Root) start() {
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
		bridge.Bridge().Call("updateId:withProtobuf:", bridge.Int64(id), bridge.Bytes(pb))

		// fmt.Println(r.root.node.debugString())
	})
}

func (r *Root) Stop() {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	if r.ticker == nil {
		return
	}
	r.ticker.Stop()
}

func (r *Root) Call(funcId string, viewId int64, args []reflect.Value) []reflect.Value {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.root.call(funcId, viewId, args)
}

// Id returns the unique identifier for r.
func (r *Root) Id() int64 {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.id
}

func (r *Root) ViewId() matcha.Id {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.root.node.id
}

// Size returns the size of r.
func (r *Root) Size() layout.Point {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	return r.size
}

// SetSize sets the size of r.
func (r *Root) SetSize(p layout.Point) {
	matcha.MainLocker.Lock()
	defer matcha.MainLocker.Unlock()

	r.size = p
	r.root.addFlag(r.root.node.id, layoutFlag)
}

type viewCacheKey struct {
	id  matcha.Id
	key string
}

// Context specifies the supporting context for building a View.
type Context struct {
	prefix string
	parent *Context

	valid     bool
	node      *node
	prevIds   map[viewCacheKey]matcha.Id
	prevNodes map[matcha.Id]*node
	skipBuild map[matcha.Id]struct{}
}

// Prev returns the view returned by the last call to Build with the given key.
func (ctx *Context) Prev(key string) View {
	return ctx.prev(key, "")
}

func (ctx *Context) prev(key string, prefix string) View {
	if ctx == nil {
		return nil
	}
	if ctx.parent != nil {
		return ctx.parent.prev(key, ctx.prefix+"|"+prefix)
	}
	if !ctx.valid {
		panic("view.Context.Prev() called on invalid context")
	}
	if ctx.node == nil {
		return nil
	}
	if prefix != "" {
		key = prefix + "|" + key
	}

	cacheKey := viewCacheKey{key: key, id: ctx.node.id}
	prevId := ctx.prevIds[cacheKey]
	prevNode := ctx.prevNodes[prevId]
	if prevNode == nil {
		return nil
	}

	v := prevNode.view
	for {
		if pv, ok := v.(*painterView); ok {
			v = pv.View
			continue
		} else if vv, ok := v.(*optionsView); ok {
			v = vv.View
			continue
		}
		break
	}

	return v
}

// PrevModel returns the last result of View.Build().
func (ctx *Context) PrevModel() *Model {
	if ctx.parent != nil {
		return ctx.PrevModel()
	}
	if !ctx.valid {
		panic("view.Context.PrevModel() called on invalid context")
	}
	if ctx.node == nil {
		return nil
	}
	return ctx.node.model
}

// NewEmbed generates a new Embed for a given key. NewEmbed is a convenience around NewEmbed(ctx.NewId(key)).
func (ctx *Context) NewEmbed(key string) Embed {
	return NewEmbed(ctx.NewId(key))
}

// NewId generates a new identifier for a given key.
func (ctx *Context) NewId(key string) matcha.Id {
	if ctx == nil {
		return matcha.Id(atomic.AddInt64(&maxId, 1))
	}
	return ctx.newId(key, "")
}

func (ctx *Context) newId(key string, prefix string) matcha.Id {
	if ctx.parent != nil {
		return ctx.parent.newId(key, ctx.prefix+"|"+prefix)
	}
	if !ctx.valid {
		panic("view.Context.Prev() called on invalid context")
	}
	if prefix != "" {
		key = prefix + "|" + key
	}

	id := matcha.Id(atomic.AddInt64(&maxId, 1))
	if ctx.node != nil {
		cacheKey := viewCacheKey{key: key, id: ctx.node.id}
		if _, ok := ctx.node.root.ids[cacheKey]; ok {
			fmt.Println("Context.NewId(): key has already been used", key)
		}
		ctx.node.root.ids[cacheKey] = id
	}
	return id
}

// // SkipBuild marks the child ids as not needing to be rebuilt.
// func (ctx *Context) SkipBuild(ids []matcha.Id) {
// 	if ctx.parent != nil {
// 		ctx.parent.SkipBuild(ids)
// 		return
// 	}

// 	if ctx.skipBuild == nil {
// 		ctx.skipBuild = map[matcha.Id]struct{}{}
// 	}
// 	for _, i := range ids {
// 		ctx.skipBuild[i] = struct{}{}
// 	}
// }

// WithPrefix returns a new Context. Calls to this Prev and NewId on this context will be prepended with key.
func (ctx *Context) WithPrefix(key string) *Context {
	return &Context{prefix: key, parent: ctx}
}

// // Id returns the identifier associated with the build context.
// func (ctx *Context) Id() matcha.Id {
// 	if ctx.parent != nil {
// 		return ctx.parent.Id()
// 	}
// 	if ctx.node == nil {
// 		return 0
// 	}
// 	return ctx.node.id
// }

func (ctx *Context) Path() []matcha.Id {
	if ctx.parent != nil {
		return ctx.parent.Path()
	}
	if ctx.node == nil {
		return []matcha.Id{0}
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

type root struct {
	node        *node
	ids         map[viewCacheKey]matcha.Id
	nodes       map[matcha.Id]*node
	middlewares []middleware

	flagMu      sync.Mutex
	updateFlags map[matcha.Id]updateFlag
}

func newRoot(v View) *root {
	id := v.Id()
	root := &root{}
	root.node = &node{
		id:   id,
		path: []matcha.Id{id},
		view: v,
		root: root,
	}
	root.updateFlags = map[matcha.Id]updateFlag{v.Id(): buildFlag}
	for _, i := range internal.Middlewares() {
		root.middlewares = append(root.middlewares, i().(middleware))
	}
	return root
}

func (root *root) addFlag(id matcha.Id, f updateFlag) {
	root.flagMu.Lock()
	defer root.flagMu.Unlock()

	root.updateFlags[id] |= f
}

func (root *root) update(size layout.Point) bool {
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
	root.updateFlags = map[matcha.Id]updateFlag{}
	return updated
}

func (root *root) MarshalProtobuf2() ([]byte, error) {
	return proto.Marshal(root.MarshalProtobuf())
}

func (root *root) MarshalProtobuf() *pb.Root {
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

func (root *root) build() {
	prevIds := root.ids
	prevNodes := root.nodes

	root.ids = map[viewCacheKey]matcha.Id{}
	root.nodes = map[matcha.Id]*node{
		root.node.id: root.node,
	}

	// Rebuild
	root.node.build(prevIds, prevNodes)

	keys := map[matcha.Id]viewCacheKey{}
	for k, v := range root.ids {
		keys[v] = k
	}
	for k, v := range prevIds {
		keys[v] = k
	}

	ids := map[viewCacheKey]matcha.Id{}
	for k := range root.nodes {
		key, ok := keys[k]
		if ok {
			ids[key] = k
		}
	}
	root.ids = ids
}

func (root *root) layout(minSize layout.Point, maxSize layout.Point) {
	g := root.node.layout(minSize, maxSize)
	g.Frame = g.Frame.Add(layout.Pt(-g.Frame.Min.X, -g.Frame.Min.Y)) // Move Frame.Min to the origin.
	root.node.layoutGuide = &g
}

func (root *root) paint() {
	root.node.paint()
}

func (root *root) call(funcId string, viewId int64, args []reflect.Value) []reflect.Value {
	node, ok := root.nodes[matcha.Id(viewId)]
	if !ok || node.model == nil {
		return nil
	}

	f, ok := node.model.NativeFuncs[funcId]
	if !ok {
		return nil
	}
	v := reflect.ValueOf(f)

	return v.Call(args)
}

type node struct {
	id    matcha.Id
	path  []matcha.Id
	root  *root
	view  View
	stage Stage

	buildId       int64
	buildPbId     int64
	buildNotify   bool
	buildNotifyId comm.Id
	model         *Model
	children      map[matcha.Id]*node

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
	for k, v := range n.model.NativeValues {
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

func (n *node) build(prevIds map[viewCacheKey]matcha.Id, prevNodes map[matcha.Id]*node) {
	if n.root.updateFlags[n.id].needsBuild() {
		n.buildId += 1

		// Send lifecycle event to new children.
		if n.stage == StageDead {
			n.view.Lifecycle(n.stage, StageVisible)
			n.stage = StageVisible
		}

		// Generate the new viewModel.
		ctx := &Context{valid: true, node: n, prevIds: prevIds, prevNodes: prevNodes}
		temp := n.view.Build(ctx)
		viewModel := &temp
		viewModelChildren := map[matcha.Id]View{}
		for _, i := range viewModel.Children {
			viewModelChildren[i.Id()] = i
		}

		// Call middleware
		for _, i := range n.root.middlewares {
			i.Build(ctx, viewModel)
		}
		ctx.valid = false

		// Diff the old children (n.children) with new children (viewModelChildren).
		addedIds := []matcha.Id{}
		removedIds := []matcha.Id{}
		unchangedIds := []matcha.Id{}
		for id := range n.children {
			if _, ok := viewModelChildren[id]; !ok {
				removedIds = append(removedIds, id)
			} else {
				unchangedIds = append(unchangedIds, id)
			}
		}
		for id := range viewModelChildren {
			if _, ok := n.children[id]; !ok {
				addedIds = append(addedIds, id)
			}
		}

		children := map[matcha.Id]*node{}
		// Add build contexts for new children.
		for _, id := range addedIds {
			var view View
			for _, i := range viewModelChildren {
				if i.Id() == id {
					view = i
					break
				}
			}

			path := make([]matcha.Id, len(n.path)+1)
			copy(path, n.path)
			path[len(n.path)] = id

			children[id] = &node{
				id:   id,
				path: path,
				view: view,
				root: n.root,
			}

			// Mark as needing rebuild
			n.root.updateFlags[id] |= buildFlag
		}
		// Mark unupdated keys as needing rebuild.
		for _, id := range unchangedIds {
			children[id] = n.children[id]
			n.root.updateFlags[id] |= buildFlag
		}
		// Send lifecycle event to removed childern.
		for _, id := range removedIds {
			n.children[id].done()
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
		i.build(prevIds, prevNodes)

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
	ctx := &layout.Context{
		MinSize:  minSize,
		MaxSize:  maxSize,
		ChildIds: []matcha.Id{},
		LayoutFunc: func(id matcha.Id, minSize, maxSize layout.Point) layout.Guide {
			child, ok := n.children[id]
			if !ok {
				fmt.Println("Attempting to layout unknown child: ", id)
				return layout.Guide{}
			}
			return child.layout(minSize, maxSize)
		},
	}
	for i := range n.children {
		ctx.ChildIds = append(ctx.ChildIds, i)
	}

	// Perform layout
	layouter := n.model.Layouter
	if layouter == nil {
		layouter = &full.Layouter{}
	}
	g, gs := layouter.Layout(ctx)
	g = g.Fit(ctx)

	// Assign guides to children
	for k, v := range gs {
		guide := v
		child, ok := n.children[k]
		if !ok {
			fmt.Println("Attempting to assign layout guide to unknown child: ", k)
			continue
		}
		child.layoutGuide = &guide
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

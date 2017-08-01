/*
Package constraint implements a constraint-based layout system.

 func (v *View) Build(ctx *view.Context) view.Model {
 	 // Create a new constraint system.
	 l := &constraint.Layouter{}

	 // Solves for the position of v, given the constraints on s. The result is a 400x100 frame.
	 l.Solve(func(s *constraint.Solver) {
	 	s.Width(400)
	 	s.Width(200) // If two constraints conflict, the earlier one is preferred.
	 	s.Height(100)
 	 }

	 // Adds a child view and solves for its position relative to v. The result is a 5x10 frame pinned to the lower right corner of v.
	 child1 := basicview.New(ctx, "child1")
	 guide1 := l.Add(child1, func(s *constraint.Solver) {
		 s.Width(5) // Left(), Top(), CenterX()... methods support constraining to floats.
		 s.Height(10)
		 s.TopEqual(l.Bottom()) // LeftEqual(), TopLess(), CenterXGreater()... methods support constraining to anchors.
		 s.LeftEqual(l.Right())
	 })

	 // Anchors can be manipulated outside of the solver function.
	 verticalCenter := l.CenterX().Add(10)

	 // Adds a child view that is twice as large as child1 and 10 points above the center v.
	 child2 := basicview.New(ctx, "child2")
	 _ = l.Add(child1, func(s *constraint.Solver) {
		 s.WidthEqual(guide1.Width().Mul(2)) // Anchors can be added to and multiplied by constants.
		 s.HeightEqual(guide1.Height().Mul(2))
		 s.CenterXEqual(l.CenterX())
		 s.CenterYEqual(verticalCenter.Add(10))
	 })

	 // Recalulates the constraints for child1.
	 guide1.Solve(func(s *constraint.Solver) {
	 	s.Width(40)
	 	s.Height(30)
	 	s.TopEqual(l.Bottom()) // The top and left position must be respecified, even though only the width and height have been updated.
	 	s.LeftEqual(l.Right())
 	 })

 	 // Solvers do not run simultaneously! Child2 is still 10x20 since at the time it was added Child1 was 5x10.

	 return view.Model{
	 	Views: l.Views(),
	 	Layouter:l,
 	 }
 }

If a child view is unconstrained in x or y, it will try to move as close to the center of the parent as possible.
If the view is unconstrained in width or height, it will try to match the minGuide as close as possible.
*/
package constraint

import (
	"fmt"
	"math"

	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal/device"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/view"
)

type comparison int

const (
	equal comparison = iota
	greater
	less
)

func (c comparison) String() string {
	switch c {
	case equal:
		return "="
	case greater:
		return ">"
	case less:
		return "<"
	}
	return ""
}

type attribute int

const (
	leftAttr attribute = iota
	rightAttr
	topAttr
	bottomAttr
	widthAttr
	heightAttr
	centerXAttr
	centerYAttr
)

func (a attribute) String() string {
	switch a {
	case leftAttr:
		return "Left"
	case rightAttr:
		return "Right"
	case topAttr:
		return "Top"
	case bottomAttr:
		return "Bottom"
	case widthAttr:
		return "Width"
	case heightAttr:
		return "Height"
	case centerXAttr:
		return "CenterX"
	case centerYAttr:
		return "CenterY"
	}
	return ""
}

// Anchor represents a float64 value that is materialized during the layout phase.
type Anchor struct {
	anchor anchor
}

// Add returns a new Anchor that is offset by v.
func (a *Anchor) Add(v float64) *Anchor {
	return &Anchor{
		offsetAnchor{
			offset:     v,
			underlying: a.anchor,
		},
	}
}

// Mul returns a new anchor that is multiplied by v.
func (a *Anchor) Mul(v float64) *Anchor {
	return &Anchor{
		multiplierAnchor{
			multiplier: v,
			underlying: a.anchor,
		},
	}
}

type anchor interface {
	value(*Layouter) float64
}

type multiplierAnchor struct {
	multiplier float64
	underlying anchor
}

func (a multiplierAnchor) value(sys *Layouter) float64 {
	return a.underlying.value(sys) * a.multiplier
}

type offsetAnchor struct {
	offset     float64
	underlying anchor
}

func (a offsetAnchor) value(sys *Layouter) float64 {
	return a.underlying.value(sys) + a.offset
}

type constAnchor float64

func (a constAnchor) value(sys *Layouter) float64 {
	return float64(a)
}

type notifierAnchor struct {
	n comm.Float64Notifier
}

func (a notifierAnchor) value(sys *Layouter) float64 {
	return a.n.Value()
}

type guideAnchor struct {
	guide     *Guide
	attribute attribute
}

func (a guideAnchor) value(sys *Layouter) float64 {
	var g layout.Guide
	switch a.guide.id {
	case rootId:
		g = *sys.Guide.matchaGuide
	case minId:
		g = *sys.min.matchaGuide
	case maxId:
		g = *sys.max.matchaGuide
	default:
		g = *sys.children[a.guide.id].matchaGuide
	}

	// if g == nil {
	// 	return 0
	// }

	switch a.attribute {
	case leftAttr:
		return g.Left()
	case rightAttr:
		return g.Right()
	case topAttr:
		return g.Top()
	case bottomAttr:
		return g.Bottom()
	case widthAttr:
		return g.Width()
	case heightAttr:
		return g.Height()
	case centerXAttr:
		return g.CenterX()
	case centerYAttr:
		return g.CenterY()
	}
	return 0
}

// Const returns a new Anchor with a constant value f.
func Const(f float64) *Anchor {
	return &Anchor{constAnchor(f)}
}

// Notifier returns a new Anchor whose value is equal to n.Value().
func Notifier(n comm.Float64Notifier) *Anchor {
	return &Anchor{notifierAnchor{n}}
}

// Guide represents a layout.Guide that is materialized during the layout phase.
type Guide struct {
	id          matcha.Id
	system      *Layouter
	children    map[matcha.Id]*Guide
	matchaGuide *layout.Guide
}

// Top returns the minimum Y coordinate as an Anchor.
func (g *Guide) Top() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: topAttr}}
}

// Right returns the maximum X coordinate as an Anchor.
func (g *Guide) Right() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: rightAttr}}
}

// Bottom returns the maximum Y coordinate as an Anchor.
func (g *Guide) Bottom() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: bottomAttr}}
}

// Left returns the minimum X coordinate as an Anchor.
func (g *Guide) Left() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: leftAttr}}
}

// Width returns the width of g as an Anchor.
func (g *Guide) Width() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: widthAttr}}
}

// Height returns the height of g as an Anchor.
func (g *Guide) Height() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: heightAttr}}
}

// CenterX returns the center of g along the X axis as an Anchor.
func (g *Guide) CenterX() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: centerXAttr}}
}

// CenterY returns the center of g along the Y axis as an Anchor.
func (g *Guide) CenterY() *Anchor {
	return &Anchor{guideAnchor{guide: g, attribute: centerYAttr}}
}

// Solve immediately calls solveFunc to update the constraints for g.
func (g *Guide) Solve(solveFunc func(*Solver)) {
	s := &Solver{id: g.id}
	if solveFunc != nil {
		solveFunc(s)
	}
	g.system.solvers = append(g.system.solvers, s)

	// Add any new notifier anchors to our notifier list.
	for _, i := range s.constraints {
		if anchor, ok := i.anchor.(notifierAnchor); ok {
			g.system.notifiers = append(g.system.notifiers, anchor.n)
		}
	}
}

func (g *Guide) add(v view.View, solveFunc func(*Solver)) *Guide {
	id := v.Id()
	chl := &Guide{
		id:          id,
		system:      g.system,
		children:    map[matcha.Id]*Guide{},
		matchaGuide: nil,
	}
	s := &Solver{id: id}
	if solveFunc != nil {
		solveFunc(s)
	}
	g.children[id] = chl
	g.system.solvers = append(g.system.solvers, s)
	g.system.views = append(g.system.views, v)

	// Add any new notifier anchors to our notifier list.
	for _, i := range s.constraints {
		if anchor, ok := i.anchor.(notifierAnchor); ok {
			g.system.notifiers = append(g.system.notifiers, anchor.n)
		}
	}
	return chl
}

type constraint struct {
	attribute  attribute
	comparison comparison
	anchor     anchor
}

func (c constraint) String() string {
	return fmt.Sprintf("%v%v%v", c.attribute, c.comparison, c.anchor)
}

// Solver is a list of constraints to be applied to a view.
type Solver struct {
	debug       bool
	id          matcha.Id
	constraints []constraint
}

func (s *Solver) solve(sys *Layouter, ctx *layout.Context) {
	cr := newConstrainedRect()

	for _, i := range s.constraints {
		copy := cr

		// Generate the range from constraint
		var r _range
		switch i.comparison {
		case equal:
			r = _range{min: i.anchor.value(sys), max: i.anchor.value(sys)}
		case greater:
			r = _range{min: i.anchor.value(sys), max: math.Inf(1)}
		case less:
			r = _range{min: math.Inf(-1), max: i.anchor.value(sys)}
		}

		// Update the solver
		switch i.attribute {
		case leftAttr:
			copy.left = copy.left.intersect(r)
		case rightAttr:
			copy.right = copy.right.intersect(r)
		case topAttr:
			copy.top = copy.top.intersect(r)
		case bottomAttr:
			copy.bottom = copy.bottom.intersect(r)
		case widthAttr:
			copy.width = copy.width.intersect(r)
		case heightAttr:
			copy.height = copy.height.intersect(r)
		case centerXAttr:
			copy.centerX = copy.centerX.intersect(r)
		case centerYAttr:
			copy.centerY = copy.centerY.intersect(r)
		}

		// Validate that the new system is well-formed. Otherwise ignore the changes.
		if !copy.isValid() {
			if s.debug {
				fmt.Println("constraint: Debug 0", i, copy) // TODO(KD): Better debugging.
			}
			continue
		}
		cr = copy
	}
	if s.debug {
		fmt.Println("constraint: Debug 1", cr, s.constraints)
	}

	// Get parent guide.
	var parent layout.Guide
	if s.id == rootId {
		parent = *sys.min.matchaGuide
	} else {
		parent = *sys.Guide.matchaGuide
	}

	// Solve for width & height.
	var width, height float64
	var g layout.Guide
	if s.id == rootId {
		g = layout.Guide{}
		width, _ = cr.solveWidth(parent.Width())
		height, _ = cr.solveHeight(parent.Height())
	} else {
		// Update the width and height ranges based on other constraints.
		_, cr = cr.solveWidth(0)
		_, cr = cr.solveHeight(0)

		g = ctx.LayoutChild(s.id, layout.Pt(cr.width.min, cr.height.min), layout.Pt(cr.width.max, cr.height.max))
		width = g.Width()
		height = g.Height()

		// Round width and height to screen scale
		width = math.Floor(width*device.ScreenScale+0.5) / device.ScreenScale

		if width < cr.width.min || height < cr.height.min || width > cr.width.max || height > cr.height.max {
			// fmt.Printf("constraint: child guide is outside of bounds. Min:%v Max:%v Actual:%v\n", layout.Pt(cr.width.min, cr.height.min), layout.Pt(cr.width.max, cr.height.max), layout.Pt(width, height))
			width = cr.width.min
			height = cr.height.min
		}
	}

	// Solve for centerX & centerY using new width & height.
	cr.width = cr.width.intersect(_range{min: width, max: width})
	cr.height = cr.height.intersect(_range{min: height, max: height})
	if !cr.isValid() {
		panic("constraint: system inconsistency")
	}
	var centerX, centerY float64
	if s.id == rootId {
		centerX = width / 2
		centerY = height / 2
	} else {
		centerX, _ = cr.solveCenterX(parent.CenterX())
		centerY, _ = cr.solveCenterY(parent.CenterY())
	}

	// Set zIndex
	g.ZIndex = sys.zIndex
	sys.zIndex += 1

	// Update the guide and the system.
	g.Frame = layout.Rt(centerX-width/2, centerY-height/2, centerX+width/2, centerY+height/2)
	if s.id == rootId {
		sys.Guide.matchaGuide = &g
	} else {
		sys.Guide.children[s.id].matchaGuide = &g
	}
	if s.debug {
		fmt.Println("constraint: Debug 2", g)
	}
}

// Debug adds debug logging for the solver.
func (s *Solver) Debug() {
	s.debug = true
}

func (s *Solver) Top(v float64) {
	s.TopEqual(Const(v))
}

func (s *Solver) TopEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: topAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) TopLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: topAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) TopGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: topAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) Right(v float64) {
	s.RightEqual(Const(v))
}

func (s *Solver) RightEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: rightAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) RightLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: rightAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) RightGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: rightAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) Bottom(v float64) {
	s.BottomEqual(Const(v))
}

func (s *Solver) BottomEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: bottomAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) BottomLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: bottomAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) BottomGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: bottomAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) Left(v float64) {
	s.LeftEqual(Const(v))
}

func (s *Solver) LeftEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: leftAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) LeftLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: leftAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) LeftGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: leftAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) Width(v float64) {
	s.WidthEqual(Const(v))
}

func (s *Solver) WidthEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: widthAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) WidthLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: widthAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) WidthGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: widthAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) Height(v float64) {
	s.HeightEqual(Const(v))
}

func (s *Solver) HeightEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: heightAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) HeightLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: heightAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) HeightGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: heightAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) CenterX(v float64) {
	s.CenterXEqual(Const(v))
}

func (s *Solver) CenterXEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerXAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) CenterXLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerXAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) CenterXGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerXAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) CenterY(v float64) {
	s.CenterYEqual(Const(v))
}

func (s *Solver) CenterYEqual(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerYAttr, comparison: equal, anchor: a.anchor})
}

func (s *Solver) CenterYLess(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerYAttr, comparison: less, anchor: a.anchor})
}

func (s *Solver) CenterYGreater(a *Anchor) {
	s.constraints = append(s.constraints, constraint{attribute: centerYAttr, comparison: greater, anchor: a.anchor})
}

func (s *Solver) String() string {
	return fmt.Sprintf("Solver{%v, %v}", s.id, s.constraints)
}

type systemId int

const (
	rootId matcha.Id = -1 * iota
	minId
	maxId
)

type Layouter struct {
	// Guide represents the size of the view that the layouter is attached to. By default Guide is the same size as MinGuide.
	Guide
	min            Guide
	max            Guide
	solvers        []*Solver
	zIndex         int
	notifiers      []comm.Notifier
	groupNotifiers map[comm.Id]notifier
	maxId          comm.Id
	views          []view.View
}

func (l *Layouter) initialize() {
	if l.groupNotifiers == nil {
		l.Guide = Guide{id: rootId, system: l, children: map[matcha.Id]*Guide{}}
		l.min = Guide{id: minId, system: l, children: map[matcha.Id]*Guide{}}
		l.max = Guide{id: maxId, system: l, children: map[matcha.Id]*Guide{}}
		l.groupNotifiers = map[comm.Id]notifier{}
	}
}

// View returns a list of all views added to l.
func (l *Layouter) Views() []view.View {
	return l.views
}

// MinGuide returns a guide representing the smallest allowed size for the view.
func (l *Layouter) MinGuide() *Guide {
	l.initialize()
	return &l.min
}

// MaxGuide returns a guide representing the largest allowed size for the view.
func (l *Layouter) MaxGuide() *Guide {
	l.initialize()
	return &l.max
}

// Layout evaluates the constraints and returns the calculated guide and child guides.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	l.initialize()
	l.min.matchaGuide = &layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
	}
	l.max.matchaGuide = &layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MaxSize.X, ctx.MaxSize.Y),
	}
	l.Guide.matchaGuide = &layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
	}
	// TODO(KD): reset all guides

	for _, i := range l.solvers {
		i.solve(l, ctx)
	}

	g := *l.Guide.matchaGuide
	gs := map[matcha.Id]layout.Guide{}
	for k, v := range l.Guide.children {
		gs[k] = *v.matchaGuide
	}
	return g, gs
}

// Add immediately calls solveFunc to generate the constraints for v. These constraints are solved by l during the layout phase.
// A corresponding guide is returned, which can be used to position other views or reposition v. If the view is not fully constrained
// it will try to match the MinGuide in dimension and center. If the child view is not fully constrained it will try to match the parent in center.
func (l *Layouter) Add(v view.View, solveFunc func(*Solver)) *Guide {
	l.initialize()
	return l.Guide.add(v, solveFunc)
}

func (l *Layouter) Solve(solveFunc func(*Solver)) {
	l.initialize()
	l.Guide.Solve(solveFunc)
}

type notifier struct {
	notifier *comm.Relay
	id       comm.Id
}

// Notify calls f anytime a Notifier anchor changes. Other updates to the layouter, such as adding guides will not trigger the notification.
func (l *Layouter) Notify(f func()) comm.Id {
	if len(l.notifiers) == 0 {
		return 0
	}

	n := &comm.Relay{}
	for _, i := range l.notifiers {
		n.Subscribe(i)
	}

	l.maxId += 1
	l.groupNotifiers[l.maxId] = notifier{
		notifier: n,
		id:       n.Notify(f),
	}
	return l.maxId
}

// Unnotify stops notifications for id.
func (l *Layouter) Unnotify(id comm.Id) {
	n, ok := l.groupNotifiers[id]
	if ok {
		n.notifier.Unnotify(n.id)
		delete(l.groupNotifiers, id)
	}
}

type _range struct {
	min float64
	max float64
}

func (r _range) intersectMin(v float64) _range {
	r.min = math.Max(r.min, v)
	return r
}

func (r _range) intersectMax(v float64) _range {
	r.max = math.Min(r.max, v)
	return r
}

func (r _range) intersect(r2 _range) _range {
	return _range{min: math.Max(r.min, r2.min), max: math.Min(r.max, r2.max)}
}

func (r _range) isValid() bool {
	if r.max < r.min {
		fmt.Println("invalid2", r.max-r.min)
	}
	return r.max >= r.min
}

func (r _range) nearest(v float64) float64 {
	// return a sane value even if range is invalid
	if r.max < r.min {
		r.max, r.min = r.min, r.max
	}
	switch {
	case r.min == r.max:
		return r.min
	case r.min >= v:
		return r.min
	case r.max <= v:
		return r.max
	default:
		return v
	}
}

type constrainedRect struct {
	left, right, top, bottom, width, height, centerX, centerY _range
}

func newConstrainedRect() constrainedRect {
	all := _range{min: math.Inf(-1), max: math.Inf(1)}
	pos := _range{min: 0, max: math.Inf(1)}
	return constrainedRect{
		left: all, right: all, top: all, bottom: all, width: pos, height: pos, centerX: all, centerY: all,
	}
}

func (cr constrainedRect) isValid() bool {
	_, r1 := cr.solveWidth(0)
	_, r2 := cr.solveHeight(0)
	_, r3 := cr.solveCenterX(0)
	_, r4 := cr.solveCenterY(0)
	return r1.width.isValid() && r2.height.isValid() && r3.centerX.isValid() && r4.centerY.isValid()
}

func (r constrainedRect) solveWidth(b float64) (float64, constrainedRect) {
	centerXMax, centerXMin := r.centerX.max, r.centerX.min
	rightMax, rightMin := r.right.max, r.right.min
	leftMax, leftMin := r.left.max, r.left.min

	// Width = (Right - CenterX) * 2
	if !math.IsInf(centerXMin, 0) && !math.IsInf(rightMax, 0) {
		r.width = r.width.intersectMax((rightMax - centerXMin) * 2)
	}
	if !math.IsInf(centerXMax, 0) && !math.IsInf(rightMin, 0) {
		r.width = r.width.intersectMin((rightMin - centerXMax) * 2)
	}

	// Width = Right - Left
	if !math.IsInf(rightMax, 0) && !math.IsInf(leftMin, 0) {
		r.width = r.width.intersectMax(rightMax - leftMin)
	}
	if !math.IsInf(rightMin, 0) && !math.IsInf(leftMax, 0) {
		r.width = r.width.intersectMin(rightMin - leftMax)
	}

	// Width = (CenterX - Left) * 2
	if !math.IsInf(centerXMax, 0) && !math.IsInf(leftMin, 0) {
		r.width = r.width.intersectMax((centerXMax - leftMin) * 2)
	}
	if !math.IsInf(centerXMin, 0) && !math.IsInf(leftMax, 0) {
		r.width = r.width.intersectMin((centerXMin - leftMax) * 2)
	}

	return r.width.nearest(b), r
}

func (r constrainedRect) solveCenterX(b float64) (float64, constrainedRect) {
	rightMax, rightMin := r.right.max, r.right.min
	leftMax, leftMin := r.left.max, r.left.min
	widthMax, widthMin := r.width.max, r.width.min

	// CenterX = (Right + Left)/2
	if !math.IsInf(rightMax, 0) && !math.IsInf(leftMax, 0) {
		r.centerX = r.centerX.intersectMax((rightMax + leftMax) / 2)
	}
	if !math.IsInf(rightMin, 0) && !math.IsInf(leftMin, 0) {
		r.centerX = r.centerX.intersectMin((rightMin + leftMin) / 2)
	}

	// CenterX = Right - Width / 2
	if !math.IsInf(rightMax, 0) && !math.IsInf(widthMin, 0) {
		r.centerX = r.centerX.intersectMax(rightMax - widthMin/2)
	}
	if !math.IsInf(rightMin, 0) && !math.IsInf(widthMax, 0) {
		r.centerX = r.centerX.intersectMin(rightMin - widthMax/2)
	}

	// CenterX = Left + Width / 2
	if !math.IsInf(leftMax, 0) && !math.IsInf(widthMax, 0) {
		r.centerX = r.centerX.intersectMax(leftMax + widthMax/2)
	}
	if !math.IsInf(leftMin, 0) && !math.IsInf(widthMin, 0) {
		r.centerX = r.centerX.intersectMin(leftMin + widthMin/2)
	}

	return r.centerX.nearest(b), r
}

func (r constrainedRect) solveHeight(b float64) (float64, constrainedRect) {
	centerYMax, centerYMin := r.centerY.max, r.centerY.min
	bottomMax, bottomMin := r.bottom.max, r.bottom.min
	topMax, topMin := r.top.max, r.top.min

	// height = (bottom - centerY) * 2
	if !math.IsInf(centerYMin, 0) && !math.IsInf(bottomMax, 0) {
		r.height = r.height.intersectMax((bottomMax - centerYMin) * 2)
	}
	if !math.IsInf(centerYMax, 0) && !math.IsInf(bottomMin, 0) {
		r.height = r.height.intersectMin((bottomMin - centerYMax) * 2)
	}

	// height = bottom - top
	if !math.IsInf(bottomMax, 0) && !math.IsInf(topMin, 0) {
		r.height = r.height.intersectMax(bottomMax - topMin)
	}
	if !math.IsInf(bottomMin, 0) && !math.IsInf(topMax, 0) {
		r.height = r.height.intersectMin(bottomMin - topMax)
	}

	// height = (centerY - top) * 2
	if !math.IsInf(centerYMax, 0) && !math.IsInf(topMin, 0) {
		r.height = r.height.intersectMax((centerYMax - topMin) * 2)
	}
	if !math.IsInf(centerYMin, 0) && !math.IsInf(topMax, 0) {
		r.height = r.height.intersectMin((centerYMin - topMax) * 2)
	}

	return r.height.nearest(b), r
}

func (r constrainedRect) solveCenterY(b float64) (float64, constrainedRect) {
	bottomMax, bottomMin := r.bottom.max, r.bottom.min
	topMax, topMin := r.top.max, r.top.min
	heightMax, heightMin := r.height.max, r.height.min

	// centerY = (bottom + top)/2
	if !math.IsInf(bottomMax, 0) && !math.IsInf(topMax, 0) {
		r.centerY = r.centerY.intersectMax((bottomMax + topMax) / 2)
	}
	if !math.IsInf(bottomMin, 0) && !math.IsInf(topMin, 0) {
		r.centerY = r.centerY.intersectMin((bottomMin + topMin) / 2)
	}

	// centerY = bottom - height / 2
	if !math.IsInf(bottomMax, 0) && !math.IsInf(heightMin, 0) {
		r.centerY = r.centerY.intersectMax(bottomMax - heightMin/2)
	}
	if !math.IsInf(bottomMin, 0) && !math.IsInf(heightMax, 0) {
		r.centerY = r.centerY.intersectMin(bottomMin - heightMax/2)
	}

	// centerY = top + height / 2
	if !math.IsInf(topMax, 0) && !math.IsInf(heightMax, 0) {
		r.centerY = r.centerY.intersectMax(topMax + heightMax/2)
	}
	if !math.IsInf(topMin, 0) && !math.IsInf(heightMin, 0) {
		r.centerY = r.centerY.intersectMin(topMin + heightMin/2)
	}

	return r.centerY.nearest(b), r
}

func (r constrainedRect) String() string {
	return fmt.Sprintf("{left:%v, right:%v, top:%v, bottom:%v, width:%v, height:%v, centerX:%v, centerY:%v}", r.left, r.right, r.top, r.bottom, r.width, r.height, r.centerX, r.centerY)
}

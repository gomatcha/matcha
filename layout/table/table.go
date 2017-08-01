/*
Package table implements a vertical, single column layout system. Views are layed out from top to bottom.

 l := &table.Layouter{}

 childView := NewChildView(...)
 l.Add(childView, nil) // The height of the view is determined by the child's layouter.

 return view.Model{
 	Views: l.Views(),
 	Layouter:l,
 }
*/
package table

import (
	"math"

	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/view"
)

type ScrollBehavior interface {
}

type Layouter struct {
	ids   []matcha.Id
	views []view.View
}

// Views returns all views that have been added to l.
func (l *Layouter) Views() []view.View {
	return l.views
}

// Add adds v to the layouter and positions it with g.
func (l *Layouter) Add(v view.View, b ScrollBehavior) {
	l.ids = append(l.ids, v.Id())
	l.views = append(l.views, v)
}

// Layout implements the view.Layouter interface.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	g := layout.Guide{}
	gs := map[matcha.Id]layout.Guide{}
	y := 0.0
	x := ctx.MinSize.X
	for i, id := range l.ids {
		g := ctx.LayoutChild(id, layout.Pt(x, 0), layout.Pt(x, math.Inf(1)))
		g.Frame = layout.Rt(0, y, g.Width(), y+g.Height())
		g.ZIndex = i
		gs[id] = g
		y += g.Height()
	}
	g.Frame = layout.Rt(0, 0, x, y)
	return g, gs
}

// Notify implements the view.Layouter interface.
func (l *Layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *Layouter) Unnotify(id comm.Id) {
	// no-op
}

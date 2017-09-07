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

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/view"
)

// Direction is the axis on which the Layouter layouts.
type Direction int

const (
	DirectionDown Direction = iota
	DirectionUp
	DirectionLeft
	DirectionRight
)

// Behavior does nothing at the moment.
type Behavior interface {
}

type Layouter struct {
	Direction Direction
	views     []view.View
}

// Views returns all views that have been added to l.
func (l *Layouter) Views() []view.View {
	return l.views
}

// Add adds v to the layouter and positions it with g. Pass nil for the behavior.
func (l *Layouter) Add(v view.View, b Behavior) {
	l.views = append(l.views, v)
}

// Layout implements the view.Layouter interface.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{}
	gs := []layout.Guide{}

	if l.Direction == DirectionDown || l.Direction == DirectionUp {
		y := 0.0
		x := ctx.MinSize.X
		for i := range l.views {
			if l.Direction == DirectionUp {
				i = len(l.views) - i - 1
			}
			g := ctx.LayoutChild(i, layout.Pt(x, 0), layout.Pt(x, math.Inf(1)))
			g.Frame = layout.Rt(0, y, g.Width(), y+g.Height())
			g.ZIndex = i
			gs = append(gs, g)
			y += g.Height()
		}
		g.Frame = layout.Rt(0, 0, x, y)
	}

	// reverse slice
	if l.Direction == DirectionUp {
		for i := len(gs)/2 - 1; i >= 0; i-- {
			opp := len(gs) - 1 - i
			gs[i], gs[opp] = gs[opp], gs[i]
		}
	}
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

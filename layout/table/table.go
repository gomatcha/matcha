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
	"fmt"
	"math"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/view"
)

// TODO(KD): Behavior does nothing at the moment.
type Behavior interface {
}

type Layouter struct {
	StartEdge layout.Edge // If no edges or more than one edge is specified layout.EdgeTop will be used.
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
func (l *Layouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{}
	gs := []layout.Guide{}

	startEdge := l.StartEdge
	switch startEdge {
	case layout.EdgeTop, layout.EdgeBottom, layout.EdgeRight, layout.EdgeLeft:
		// no-op
	default:
		startEdge = layout.EdgeTop
	}

	if startEdge == layout.EdgeBottom || startEdge == layout.EdgeTop {
		y := 0.0
		x := ctx.MinSize().X
		if x == 0 {
			fmt.Println("table.Layouter: Width is 0, is scrollview.ScrollAxes set?")
		}
		for i := range l.views {
			if startEdge == layout.EdgeBottom {
				i = len(l.views) - i - 1
			}
			g := ctx.LayoutChild(i, layout.Pt(x, 0), layout.Pt(x, math.Inf(1)))
			g.Frame = layout.Rt(0, y, g.Width(), y+g.Height())
			g.ZIndex = i
			gs = append(gs, g)
			y += g.Height()
		}
		g.Frame = layout.Rt(0, 0, x, y)
	} else {
		y := ctx.MinSize().Y
		x := 0.0
		if y == 0 {
			fmt.Println("table.Layouter: Height is 0, is scrollview.ScrollAxes set?")
		}
		for i := range l.views {
			if startEdge == layout.EdgeLeft {
				i = len(l.views) - i - 1
			}
			g := ctx.LayoutChild(i, layout.Pt(0, y), layout.Pt(math.Inf(1), y))
			g.Frame = layout.Rt(x, 0, x+g.Width(), g.Height())
			g.ZIndex = i
			gs = append(gs, g)
			x += g.Width()
		}
		g.Frame = layout.Rt(0, 0, x, y)
	}

	// reverse slice
	if startEdge == layout.EdgeBottom || startEdge == layout.EdgeLeft {
		for i := len(gs)/2 - 1; i >= 0; i-- {
			opp := len(gs) - 1 - i
			gs[i], gs[opp] = gs[opp], gs[i]
		}
	}
	return g, gs
}

// DebugStrings must be called after Layout()...
func (l *Layouter) DebugStrings() (string, []string) {
	debugstr := "Edge="
	switch l.StartEdge {
	case layout.EdgeBottom:
		debugstr += "Bot"
	case layout.EdgeLeft:
		debugstr += "Left"
	case layout.EdgeRight:
		debugstr += "Right"
	case layout.EdgeTop:
	default:
		debugstr += "Top"
	}

	debugstrs := make([]string, len(l.views))
	for i := range l.views {
		debugstrs[i] = debugstr
	}
	return debugstr, debugstrs
}

// Notify implements the view.Layouter interface.
func (l *Layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *Layouter) Unnotify(id comm.Id) {
	// no-op
}

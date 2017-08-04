/*
Package layout provides geometric primitives and interfaces for view layout.

Layouter

While view.View handles the rendering of a component, the actual positioning of the
views is delegated to the Layouter interface. Each view can specify it's Layouter
in the view.Model returned by Build(). If no layouter is given, all of
its children will be positioned to size of its parent view.

Understanding the details of this is not too important for most day to day development.
For the most part you will be using predefined Layouters such as the one provided by
the constraint package.

Layout occurs in a separate pass after the view hierarchy has been built. Like view.Build(),
each layouter is only responsible for returning its own frame and the frame of its direct descendents.
To determine the correct sizing for a child, the layouter will call *Context.LayoutChild()
passing a minSize and a maxSize. The child will return a desired size within the
min and max, which the parent can then position. Here is an example Layout function
that centers its children within itself.

	func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
		// Specify that the view wants to be the minSize given by its parent.
		g := layout.Guide{
			Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
		}

		// Iterate over all child ids.
		gs := []layout.Guide{}
		for i := 0; i< ctx.ChildCount; i++ {

			// Get the desired size of the children. In this case we let the children be any size.
			child := ctx.LayoutChild(idx, layout.Pt(0, 0), layout.Pt(math.Inf(1), math.Inf(1)))

			// Position the children to be centered in the view.
			child.Frame = child.Frame.Add(layout.Pt(g.CenterX()-child.Width()/2, g.CenterY()-child.Height()/2))
			child.ZIndex = i
			gs = append(gs, child)
		}

		// Return the view's size, and the frames of its children.
		return g, gs
	}

Layouters also implement the comm.Notifier interface. This allows layouts to update
without rebuilding the view. It is light-weight and useful for animations.
*/
package layout

import (
	"reflect"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/comm"
	pblayout "gomatcha.io/matcha/pb/layout"
)

func init() {
	bridge.RegisterType("layout.Point", reflect.TypeOf(Point{}))
	bridge.RegisterType("layout.Rect", reflect.TypeOf(Rect{}))
}

type Layouter interface {
	Layout(ctx *Context) (Guide, []Guide)
	comm.Notifier
}

type Context struct {
	MinSize    Point
	MaxSize    Point
	ChildCount int
	LayoutFunc func(int, Point, Point) Guide
}

func (l *Context) LayoutChild(idx int, minSize, maxSize Point) Guide {
	g := l.LayoutFunc(idx, minSize, maxSize)
	g.Frame = g.Frame.Add(Pt(-g.Frame.Min.X, -g.Frame.Min.Y))
	return g
}

// Guide represents the position of a view.
type Guide struct {
	Frame  Rect
	ZIndex int
}

// MarshalProtobuf serializes g into a protobuf object.
func (g Guide) MarshalProtobuf() *pblayout.Guide {
	return &pblayout.Guide{
		Frame:  g.Frame.MarshalProtobuf(),
		ZIndex: int64(g.ZIndex),
	}
}

// Left returns the left edge of g.
func (g Guide) Left() float64 {
	return g.Frame.Min.X
}

// Right returns the right edge of g.
func (g Guide) Right() float64 {
	return g.Frame.Max.X
}

// Top returns the top edge of g.
func (g Guide) Top() float64 {
	return g.Frame.Min.Y
}

// Bottom returns the bottom edge of g.
func (g Guide) Bottom() float64 {
	return g.Frame.Max.Y
}

// Width returns the width of g.
func (g Guide) Width() float64 {
	return g.Frame.Max.X - g.Frame.Min.X
}

// Height returns the height of g.
func (g Guide) Height() float64 {
	return g.Frame.Max.Y - g.Frame.Min.Y
}

// CenterX returns the horizontal center of g.
func (g Guide) CenterX() float64 {
	return (g.Frame.Max.X - g.Frame.Min.X) / 2
}

// CenterY returns the vertical center of g.
func (g Guide) CenterY() float64 {
	return (g.Frame.Max.Y - g.Frame.Min.Y) / 2
}

// Fit adjusts the frame of the guide to be within MinSize and MaxSize of the LayoutContext.
func (g Guide) Fit(ctx *Context) Guide {
	if g.Width() < ctx.MinSize.X {
		g.Frame.Max.X = ctx.MinSize.X - g.Frame.Min.X
	}
	if g.Height() < ctx.MinSize.Y {
		g.Frame.Max.Y = ctx.MinSize.Y - g.Frame.Min.Y
	}
	if g.Width() > ctx.MaxSize.X {
		g.Frame.Max.X = ctx.MaxSize.X - g.Frame.Min.X
	}
	if g.Height() > ctx.MaxSize.Y {
		g.Frame.Max.Y = ctx.MaxSize.Y - g.Frame.Min.Y
	}
	return g
}

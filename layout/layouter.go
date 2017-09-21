// Package layout provides geometric primitives and interfaces for view layout.
// See https://gomatcha.io/guide/layout for more details.
package layout

import (
	"reflect"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/comm"
	pblayout "gomatcha.io/matcha/proto/layout"
)

func init() {
	bridge.RegisterType("layout.Point", reflect.TypeOf(Point{}))
	bridge.RegisterType("layout.Rect", reflect.TypeOf(Rect{}))
}

// Layouter ... TODO(KD):...
type Layouter interface {
	Layout(ctx Context) (Guide, []Guide)
	comm.Notifier
}

// Context provides properties and hooks for Layouter.Layout().
type Context interface {
	MaxSize() Point
	MinSize() Point
	ChildCount() int
	LayoutChild(idx int, minSize, maxSize Point) Guide
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

// Package layout provides geometric primitives and interfaces for view layout.
package layout

import (
	"reflect"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	pblayout "gomatcha.io/matcha/pb/layout"
)

func init() {
	bridge.RegisterType("layout.Point", reflect.TypeOf(Point{}))
	bridge.RegisterType("layout.Rect", reflect.TypeOf(Rect{}))
}

type Layouter interface {
	Layout(ctx *Context) (Guide, map[matcha.Id]Guide)
	comm.Notifier
}

type Context struct {
	MinSize    Point
	MaxSize    Point
	ChildIds   []matcha.Id
	LayoutFunc func(matcha.Id, Point, Point) Guide
}

func (l *Context) LayoutChild(id matcha.Id, minSize, maxSize Point) Guide {
	return l.LayoutFunc(id, minSize, maxSize)
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

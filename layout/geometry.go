package layout

import (
	"fmt"

	"gomatcha.io/matcha/comm"
	pblayout "gomatcha.io/matcha/pb/layout"
)

type Axis int

const (
	AxisY Axis = 1 << iota
	AxisX
)

type Direction int

const (
	DirectionUp Direction = 1 << iota
	DirectionDown
	DirectionLeft
	DirectionRight
)

// Rect represents a 2D rectangle with the top left corner at Min and the bottom
// right corner at Max.
type Rect struct {
	Min, Max Point
}

// Rt creates a rectangle with the min at X0, Y0 and the max at X1, Y1.
func Rt(x0, y0, x1, y1 float64) Rect {
	return Rect{Min: Point{X: x0, Y: y0}, Max: Point{X: x1, Y: y1}}
}

// MarshalProtobuf serializes r into a protobuf object.
func (r *Rect) MarshalProtobuf() *pblayout.Rect {
	return &pblayout.Rect{
		Min: r.Min.MarshalProtobuf(),
		Max: r.Max.MarshalProtobuf(),
	}
}

// UnmarshalProtobuf deserializes r from a protobuf object.
func (r *Rect) UnmarshalProtobuf(pbrect *pblayout.Rect) {
	r.Min.UnmarshalProtobuf(pbrect.Min)
	r.Max.UnmarshalProtobuf(pbrect.Max)
}

// Add translates rect r by p.
func (r Rect) Add(p Point) Rect {
	n := r
	n.Min.X += p.X
	n.Min.Y += p.Y
	n.Max.X += p.X
	n.Max.Y += p.Y
	return n
}

// String returns a string description of r.
func (r Rect) String() string {
	return fmt.Sprintf("Rect{%v, %v, %v, %v}", r.Min.X, r.Min.Y, r.Max.X, r.Max.Y)
}

// Point represents a point on the XY coordinate system.
type Point struct {
	X float64
	Y float64
}

// Pt creates a point with x and y.
func Pt(x, y float64) Point {
	return Point{X: x, Y: y}
}

// MarshalProtobuf serializes p into a protobuf object.
func (p *Point) MarshalProtobuf() *pblayout.Point {
	return &pblayout.Point{
		X: p.X,
		Y: p.Y,
	}
}

// UnmarshalProtobuf deserializes p from a protobuf object.
func (p *Point) UnmarshalProtobuf(pbpoint *pblayout.Point) {
	p.X = pbpoint.X
	p.Y = pbpoint.Y
}

// PointNotifier wraps the comm.Notifier interface with an additional Value() method which returns a Point.
type PointNotifier interface {
	comm.Notifier
	Value() Point
}

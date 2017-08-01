// Package paint implements view display properties.
//  func (v *View) Build(ctx *view.Context) view.Model {
//  	return view.Model{
//  		Painter: &paint.Style{
//  			BackgroundColor: colornames.Green,
//  			CornerRadius: 3,
//  		},
//  	}
//  }
package paint

import (
	"image/color"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/pb"
	"gomatcha.io/matcha/pb/paint"
)

// Painter is the interface that describes how a view should be drawn on screen.
type Painter interface {
	PaintStyle() Style
	comm.Notifier
}

// Style is a list of display properties of that views can set.
type Style struct {
	Transparency    float64
	BackgroundColor color.Color
	BorderColor     color.Color
	BorderWidth     float64
	CornerRadius    float64
	ShadowRadius    float64
	ShadowOffset    layout.Point
	ShadowColor     color.Color
}

func (s *Style) MarshalProtobuf() *paint.Style {
	return &paint.Style{
		Transparency:    s.Transparency,
		BackgroundColor: pb.ColorEncode(s.BackgroundColor),
		BorderColor:     pb.ColorEncode(s.BorderColor),
		BorderWidth:     s.BorderWidth,
		CornerRadius:    s.CornerRadius,
		ShadowRadius:    s.ShadowRadius,
		ShadowOffset:    s.ShadowOffset.MarshalProtobuf(),
		ShadowColor:     pb.ColorEncode(s.ShadowColor),
	}
}

// PaintStyle implements the Painter interface.
func (s *Style) PaintStyle() Style {
	if s == nil {
		return Style{}
	}
	return *s
}

// Notify implements the Painter interface. This is a no-op.
func (s *Style) Notify(func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the Painter interface. This is a no-op.
func (s *Style) Unnotify(id comm.Id) {
	// no-op
}

type notifier struct {
	notifier *comm.Relay
	id       comm.Id
}

// AnimatedStyle is the animated version of Style.
type AnimatedStyle struct {
	Style           Style
	Transparency    comm.Float64Notifier
	BackgroundColor comm.ColorNotifier
	BorderColor     comm.ColorNotifier
	BorderWidth     comm.Float64Notifier
	CornerRadius    comm.Float64Notifier
	ShadowRadius    comm.Float64Notifier
	ShadowOffset    layout.PointNotifier
	ShadowColor     comm.ColorNotifier

	maxId          comm.Id
	groupNotifiers map[comm.Id]notifier
}

// PaintStyle implements the Painter interface.
func (as *AnimatedStyle) PaintStyle() Style {
	s := as.Style
	if as.Transparency != nil {
		s.Transparency = as.Transparency.Value()
	}
	if as.BackgroundColor != nil {
		s.BackgroundColor = as.BackgroundColor.Value()
	}
	if as.BorderColor != nil {
		s.BorderColor = as.BorderColor.Value()
	}
	if as.BorderWidth != nil {
		s.BorderWidth = as.BorderWidth.Value()
	}
	if as.CornerRadius != nil {
		s.CornerRadius = as.CornerRadius.Value()
	}
	if as.ShadowRadius != nil {
		s.ShadowRadius = as.ShadowRadius.Value()
	}
	if as.ShadowOffset != nil {
		s.ShadowOffset = as.ShadowOffset.Value()
	}
	if as.ShadowColor != nil {
		s.ShadowColor = as.ShadowColor.Value()
	}
	return s
}

// Notify implements the Painter interface.
func (as *AnimatedStyle) Notify(f func()) comm.Id {
	n := &comm.Relay{}

	if as.Transparency != nil {
		n.Subscribe(as.Transparency)
	}
	if as.BackgroundColor != nil {
		n.Subscribe(as.BackgroundColor)
	}
	if as.BorderColor != nil {
		n.Subscribe(as.BorderColor)
	}
	if as.BorderWidth != nil {
		n.Subscribe(as.BorderWidth)
	}
	if as.CornerRadius != nil {
		n.Subscribe(as.CornerRadius)
	}
	if as.ShadowRadius != nil {
		n.Subscribe(as.ShadowRadius)
	}
	if as.ShadowOffset != nil {
		n.Subscribe(as.ShadowOffset)
	}
	if as.ShadowColor != nil {
		n.Subscribe(as.ShadowColor)
	}

	as.maxId += 1
	if as.groupNotifiers == nil {
		as.groupNotifiers = map[comm.Id]notifier{}
	}
	as.groupNotifiers[as.maxId] = notifier{
		notifier: n,
		id:       n.Notify(f),
	}
	return as.maxId
}

// Unnotify implements the Painter interface.
func (as *AnimatedStyle) Unnotify(id comm.Id) {
	n, ok := as.groupNotifiers[id]
	if ok {
		n.notifier.Unnotify(n.id)
		delete(as.groupNotifiers, id)
	}
}

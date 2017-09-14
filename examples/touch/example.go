// Package paint provides examples of how to use the matcha/touch package.
package touch

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/touch New", func() *view.Root {
		return view.NewRoot(New())
	})
}

type TouchView struct {
	view.Embed
	tapCounter    int
	pressCounter  int
	buttonCounter int
}

func New() *TouchView {
	return &TouchView{
		Embed: view.Embed{},
	}
}

func (v *TouchView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := NewTapChildView()
	chl1.OnTouch = func() {
		fmt.Println("On Tap")
		v.tapCounter += 1
		go v.Signal() // TODO(KD): Why is this on separate thread?
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.WidthEqual(constraint.Const(100))
		s.HeightEqual(constraint.Const(100))
	})

	chl2 := view.NewTextView()
	chl2.String = fmt.Sprintf("tap: %v", v.tapCounter)
	chl2.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
	})

	chl3 := NewPressChildView()
	chl3.OnTouch = func() {
		fmt.Println("On Press")
		v.pressCounter += 1
		go v.Signal()
	}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g2.Left())
		s.WidthEqual(constraint.Const(100))
		s.HeightEqual(constraint.Const(100))
	})

	chl4 := view.NewTextView()
	chl4.String = fmt.Sprintf("Press: %v", v.pressCounter)
	chl4.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g4 := l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g3.Bottom())
		s.LeftEqual(g3.Left())
	})

	chl5 := NewButtonChildView()
	chl5.OnTouch = func() {
		fmt.Println("On Button")
		v.buttonCounter += 1
		go v.Signal()
	}
	g5 := l.Add(chl5, func(s *constraint.Solver) {
		s.TopEqual(g4.Bottom())
		s.LeftEqual(g4.Left())
		s.WidthEqual(constraint.Const(100))
		s.HeightEqual(constraint.Const(100))
	})
	chl6 := view.NewTextView()
	chl6.String = fmt.Sprintf("Button: %v", v.buttonCounter)
	chl6.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g6 := l.Add(chl6, func(s *constraint.Solver) {
		s.TopEqual(g5.Bottom())
		s.LeftEqual(g5.Left())
	})
	_ = g6

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type TapChildView struct {
	view.Embed
	OnTouch func()
}

func NewTapChildView() *TapChildView {
	return &TapChildView{}
}

func (v *TapChildView) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.GestureList{&touch.TapGesture{
				Count: 1,
				OnTouch: func(e *touch.TapEvent) {
					v.OnTouch()
				},
			}},
		},
	}
}

type PressChildView struct {
	view.Embed
	OnTouch func()
}

func NewPressChildView() *PressChildView {
	return &PressChildView{}
}

func (v *PressChildView) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.GestureList{&touch.PressGesture{
				MinDuration: time.Second / 2,
				OnTouch: func(e *touch.PressEvent) {
					v.OnTouch()
				},
			}},
		},
	}
}

type ButtonChildView struct {
	view.Embed
	OnTouch func()
}

func NewButtonChildView() *ButtonChildView {
	return &ButtonChildView{}
}

func (v *ButtonChildView) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.GestureList{&touch.ButtonGesture{
				OnTouch: func(e *touch.ButtonEvent) {
					v.OnTouch()
				},
			}},
		},
	}
}

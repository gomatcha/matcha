// Package pointer provides examples of how to use the matcha/pointer package.
package pointer

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/touch NewTouchView", func() view.View {
		return NewTouchView()
	})
}

type TouchView struct {
	view.Embed
	tapCounter    int
	pressCounter  int
	buttonCounter int
}

func NewTouchView() *TouchView {
	return &TouchView{
		Embed: view.Embed{},
	}
}

func (v *TouchView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	scrollLayouter := &constraint.Layouter{}
	scrollLayouter.Solve(func(s *constraint.Solver) {
		s.Width(200)
		s.Height(500)
	})
	tap := NewTapChildView()
	tap.OnTouch = func() {
		v.tapCounter += 1
		go v.Signal() // TODO(KD): Why is this on separate thread?
	}
	scrollLayouter.Add(tap, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.Width(200)
		s.Height(100)
	})

	scrollview := view.NewScrollView()
	scrollview.ContentChildren = scrollLayouter.Views()
	scrollview.ContentPainter = &paint.Style{BackgroundColor: colornames.Yellow}
	scrollview.ContentLayouter = scrollLayouter
	g1 := l.Add(scrollview, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(200)
		s.Height(200)
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
		s.Width(100)
		s.Height(100)
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
		s.Width(100)
		s.Height(100)
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
	l := &constraint.Layouter{}
	a := view.NewBasicView()
	a.Painter = &paint.Style{BackgroundColor: colornames.Red}
	l.Add(a, func(s *constraint.Solver) {
		s.Left(0)
		s.Top(0)
		s.Width(100)
		s.Height(100)
	})
	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			pointer.GestureList{&pointer.TapGesture{
				Count: 1,
				OnEvent: func(e *pointer.TapEvent) {
					if e.Kind == pointer.EventKindPossible {
						fmt.Println("Tap Possible")
					} else if e.Kind == pointer.EventKindChanged {
						fmt.Println("Tap Changed")
					} else if e.Kind == pointer.EventKindFailed {
						fmt.Println("Tap Failed")
					} else if e.Kind == pointer.EventKindRecognized {
						fmt.Println("Tap Recognized")
						v.OnTouch()
					}
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
			pointer.GestureList{&pointer.PressGesture{
				MinDuration: time.Second / 2,
				OnEvent: func(e *pointer.PressEvent) {
					if e.Kind == pointer.EventKindPossible {
						fmt.Println("Press Possible")
					} else if e.Kind == pointer.EventKindChanged {
						fmt.Println("Press Changed")
					} else if e.Kind == pointer.EventKindFailed {
						fmt.Println("Press Failed")
					} else if e.Kind == pointer.EventKindRecognized {
						fmt.Println("Press Recognized")
						v.OnTouch()
					}
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
			pointer.GestureList{&pointer.ButtonGesture{
				OnEvent: func(e *pointer.ButtonEvent) {
					if e.Kind == pointer.EventKindPossible {
						fmt.Println("Button Possible")
					} else if e.Kind == pointer.EventKindChanged {
						fmt.Println("Button Changed")
					} else if e.Kind == pointer.EventKindFailed {
						fmt.Println("Button Failed")
					} else if e.Kind == pointer.EventKindRecognized {
						fmt.Println("Button Recognized")
						v.OnTouch()
					}
				},
			}},
		},
	}
}

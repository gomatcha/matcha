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
	counter      int
	pressCounter int
}

func New() *TouchView {
	return &TouchView{
		Embed: view.Embed{},
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := NewTouchChildView()
	chl1.OnTouch = func() {
		fmt.Println("On touch")
		v.counter += 1
		go v.Signal()
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.WidthEqual(constraint.Const(100))
		s.HeightEqual(constraint.Const(100))
	})

	chl2 := view.NewTextView()
	chl2.String = fmt.Sprintf("Counter: %v", v.counter)
	chl2.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
	})

	chl3 := NewPressChildView()
	chl3.OnPress = func() {
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
		fmt.Println("On touch")
		v.counter += 1
		go v.Signal()
	}
	g5 := l.Add(chl5, func(s *constraint.Solver) {
		s.TopEqual(g4.Bottom())
		s.LeftEqual(g4.Left())
		s.WidthEqual(constraint.Const(100))
		s.HeightEqual(constraint.Const(100))
	})
	_ = g5

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type PressChildView struct {
	view.Embed
	OnPress func()
}

func NewPressChildView() *PressChildView {
	return &PressChildView{}
}

func (v *PressChildView) Build(ctx *view.Context) view.Model {
	tap := &touch.PressRecognizer{
		MinDuration: time.Second / 2,
		OnTouch: func(e *touch.PressEvent) {
			v.OnPress()
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
		// Values: map[string]interface{}{
		// 	touch.Key: []touch.Recognizer{tap},
		// },
		// Options: []view.Options{
		// 	touch.Recognizers([]touch.Recognizer{tap}),
		// 	app.ActivityIndicator{},
		// 	app.StatusBar{
		// 		Hidden: true,
		// 		Style:StatusBarStyleLight,
		// 	},
		// },
	}
}

type TouchChildView struct {
	view.Embed
	OnTouch func()
}

func NewTouchChildView() *TouchChildView {
	return &TouchChildView{}
}

func (v *TouchChildView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapRecognizer{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			v.OnTouch()
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.RecognizerList{tap},
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

func (v *ButtonChildView) Build(ctx *view.Context) view.Model {
	button := &touch.ButtonRecognizer{
		OnTouch: func(e *touch.ButtonEvent) {
			fmt.Println("On Touch:s", e.Kind)
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			touch.RecognizerList{button},
		},
	}
}

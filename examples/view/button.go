package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/alert"
	"gomatcha.io/matcha/view/button"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewButtonView", func() *view.Root {
		return view.NewRoot(NewButtonView())
	})
}

type ButtonView struct {
	view.Embed
	value *comm.Float64Value
}

func NewButtonView() *ButtonView {
	return &ButtonView{
		value: &comm.Float64Value{},
	}
}

func (v *ButtonView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := button.New()
	chl1.Text = "Press Me"
	chl1.OnPress = func() {
		alert.Alert("Button Pressed", "")
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := button.New()
	chl2.Text = "Press Me"
	chl2.Color = colornames.Red
	chl2.Enabled = false
	chl2.OnPress = func() {
		alert.Alert("Button2 Pressed", "")
	}
	l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(50))
		s.LeftEqual(g1.Left())
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

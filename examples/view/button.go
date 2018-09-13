package view

import (
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view NewButtonView", func() view.View {
		return NewButtonView()
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

func (v *ButtonView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.Enabled = true
	chl1.String = "Press Me"
	chl1.OnPress = func() {
		view.Alert("Button Pressed", "")
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := view.NewButton()
	chl2.String = "Press Me"
	chl2.Color = colornames.Red
	chl2.Enabled = false
	chl2.OnPress = func() {
		view.Alert("Button2 Pressed", "")
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

package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewButtonView", func() view.View {
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

	label := view.NewTextView()
	label.String = "Enabled Button:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	button := view.NewButton()
	button.Enabled = true
	button.String = "Button"
	button.OnPress = func() {
		view.Alert("Button Pressed", "")
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Disabled Button:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	chl2 := view.NewButton()
	chl2.String = "Button"
	chl2.Enabled = false
	chl2.OnPress = func() {
		view.Alert("Button Pressed", "")
	}
	l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

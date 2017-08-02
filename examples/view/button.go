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
		return view.NewRoot(NewButtonView(nil, ""))
	})
}

type ButtonView struct {
	view.Embed
	value *comm.Float64Value
}

func NewButtonView(ctx *view.Context, key string) *ButtonView {
	if v, ok := ctx.Prev(key).(*ButtonView); ok {
		return v
	}
	return &ButtonView{
		Embed: ctx.NewEmbed(key),
		value: &comm.Float64Value{},
	}
}

func (v *ButtonView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := button.New(ctx, "0")
	chl.Text = "Press Me"
	chl.OnPress = func() {
		alert.Alert("Button Pressed", "")
	}
	l.Add(chl, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

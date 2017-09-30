package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewUnknownView", func() view.View {
		return NewUnknownView()
	})
}

type UnknownView struct {
	view.Embed
}

func NewUnknownView() *UnknownView {
	return &UnknownView{}
}

func (v *UnknownView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	_ = l.Add(&UnknownTestView{}, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
		s.Height(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Blue},
	}
}

type UnknownTestView struct {
	view.Embed
}

func (v *UnknownTestView) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter:        &paint.Style{BackgroundColor: colornames.Green},
		NativeViewName: "Unknown",
	}
}

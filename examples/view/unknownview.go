package view

import (
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view NewUnknownView", func() view.View {
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

// Package custom provides examples of how to import custom components.
package customview

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/examples/custom/customview"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/customview NewView", func() view.View {
		return NewView()
	})
}

type View struct {
	view.Embed
}

func NewView() *View {
	return &View{}
}

func (v *View) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := customview.NewCustomView()
	chl1.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	l.Add(chl1, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(100)
		s.Height(100)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

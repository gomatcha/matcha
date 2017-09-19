// Package paint provides examples of how to use the matcha/paint package.
package paint

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/paint New", func() view.View {
		return NewPaintView()
	})
}

type PaintView struct {
	view.Embed
}

func NewPaintView() *PaintView {
	return &PaintView{}
}

func (v *PaintView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewBasicView()
	chl1.Painter = &paint.Style{
		BackgroundColor: colornames.Blue,
		BorderColor:     colornames.Red,
		BorderWidth:     3,
		CornerRadius:    20,
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(50))
		s.LeftEqual(constraint.Const(100))
		s.Width(100)
		s.Height(100)
	})

	chl2 := view.NewBasicView()
	chl2.Painter = &paint.Style{
		BackgroundColor: colornames.Yellow,
		ShadowRadius:    4,
		ShadowOffset:    layout.Pt(5, 5),
		ShadowColor:     colornames.Black,
	}
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(20))
		s.LeftEqual(g1.Left())
		s.Width(100)
		s.Height(100)
	})

	chl3 := view.NewBasicView()
	chl3.Painter = &paint.Style{
		BackgroundColor: colornames.Black,
		Transparency:    0.8,
	}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom().Add(20))
		s.LeftEqual(g2.Left())
		s.Width(100)
		s.Height(100)
	})

	chl4 := view.NewBasicView()
	chl4.Painter = &paint.Style{BackgroundColor: colornames.Magenta}
	_ = l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g3.Bottom().Add(20))
		s.LeftEqual(g3.Left())
		s.Width(100)
		s.Height(100)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

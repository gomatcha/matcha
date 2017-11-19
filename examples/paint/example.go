// Package paint provides examples of how to use the matcha/paint package.
package paint

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/paint NewPaintView", func() view.View {
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

	label := view.NewTextView()
	label.String = "Background Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	chl := view.NewBasicView()
	chl.Painter = &paint.Style{
		BackgroundColor: colornames.Green,
	}
	g = l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Transparency:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	chl = view.NewBasicView()
	chl.Painter = &paint.Style{
		BackgroundColor: colornames.Red,
		Transparency:    0.5,
	}
	g = l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Border Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	chl = view.NewBasicView()
	chl.Painter = &paint.Style{
		BorderColor: colornames.Red,
		BorderWidth: 3,
	}
	g = l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Shadow:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	chl = view.NewBasicView()
	chl.Painter = &paint.Style{
		BackgroundColor: colornames.Yellow,
		ShadowRadius:    4,
		ShadowOffset:    layout.Pt(5, 5),
		ShadowColor:     colornames.Black,
	}
	g = l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Corner Radius:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.LeftEqual(l.CenterX().Add(15))
	})

	chl = view.NewBasicView()
	chl.Painter = &paint.Style{
		BackgroundColor: colornames.Blue,
		BorderColor:     colornames.Pink,
		BorderWidth:     2,
		CornerRadius:    20,
	}
	g = l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

package layout

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/layout NewConstraintsView", func() view.View {
		return NewConstraintsView()
	})
}

type ConstraintsView struct {
	view.Embed
}

func NewConstraintsView() *ConstraintsView {
	return &ConstraintsView{}
}

func (v *ConstraintsView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Top Left"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
	})

	label = view.NewTextView()
	label.String = "Width Height"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Yellow}
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "2x Width Height"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Pink}
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.WidthEqual(g.Width().Mul(2))
		s.HeightEqual(g.Height().Mul(2))
	})

	label = view.NewTextView()
	label.String = "Inset"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Purple}
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Top().Add(5))
		s.RightEqual(g.Right().Add(-5))
		s.WidthLess(constraint.Const(50))
		s.HeightLess(constraint.Const(50))
	})

	label = view.NewTextView()
	label.String = "CenterX CenterY"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Green}
	g = l.Add(label, func(s *constraint.Solver) {
		s.CenterX(0)
		s.CenterY(0)
	})

	label = view.NewTextView()
	label.String = "Bottom Right"
	label.Style.SetFont(text.DefaultFont(18))
	label.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
	g = l.Add(label, func(s *constraint.Solver) {
		s.Bottom(0)
		s.Right(0)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

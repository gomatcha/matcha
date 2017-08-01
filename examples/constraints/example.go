package constraints

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/constraints New", func() *view.Root {
		return view.NewRoot(New(nil, ""))
	})
}

type ConstraintsView struct {
	view.Embed
}

func New(ctx *view.Context, key string) *ConstraintsView {
	if v, ok := ctx.Prev(key).(*ConstraintsView); ok {
		return v
	}
	return &ConstraintsView{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *ConstraintsView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := basicview.New(ctx, "1")
	chl1.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	_ = l.Add(chl1, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(100)
		s.Height(100)
	})

	chl2 := basicview.New(ctx, "2")
	chl2.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		// s.TopEqual(g1.Bottom())
		// s.LeftEqual(g1.Left())
		s.Width(300)
		s.Height(300)
	})

	chl3 := basicview.New(ctx, "3")
	chl3.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g2.Left())
		s.Width(100)
		s.Height(100)
	})

	chl4 := basicview.New(ctx, "4")
	chl4.Painter = &paint.Style{BackgroundColor: colornames.Magenta}
	_ = l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g3.Right())
		s.Width(50)
		s.Height(50)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

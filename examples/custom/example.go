package custom

import (
	"github.com/overcyn/customview"
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/custom New", func() *view.Root {
		return view.NewRoot(New(nil, ""))
	})
}

type View struct {
	view.Embed
}

func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := customview.New(ctx, "1")
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

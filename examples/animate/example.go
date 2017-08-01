package animate

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/animate New", func() *view.Root {
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

func (v *View) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageVisible) {
		time.AfterFunc(time.Second*2, func() {
			fmt.Println("Update")
			v.Signal()
		})
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := basicview.New(ctx, "")
	// chl.Painter = &paint.AnimatedStyle{BackgroundColor: v.colorTicker}
	l.Add(chl, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		// s.WidthEqual(constraint.Notifier(v.floatTicker))
		// s.HeightEqual(constraint.Notifier(v.floatTicker))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

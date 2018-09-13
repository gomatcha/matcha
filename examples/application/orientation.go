package application

import (
	"strconv"

	"github.com/gomatcha/matcha/application"
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/application NewOrientationView", func() view.View {
		return NewOrientationView()
	})
}

type OrientationView struct {
	view.Embed
	notifier comm.IntNotifier
}

func NewOrientationView() *OrientationView {
	return &OrientationView{
		notifier: application.OrientationNotifier,
	}
}

func (v *OrientationView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.notifier)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.notifier)
	}
}

func (v *OrientationView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.String = "Get orientation"
	chl1.OnPress = func() {
		o := application.Orientation()
		view.Alert("Orientation"+strconv.Itoa(int(o)), "")
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := view.NewTextView()
	chl2.String = "Orientation" + strconv.Itoa(v.notifier.Value())
	_ = l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(50))
		s.LeftEqual(g1.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

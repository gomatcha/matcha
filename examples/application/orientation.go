package application

import (
	"strconv"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/application NewOrientationView", func() view.View {
		return NewOrientationView()
	})
}

type OrientationView struct {
	view.Embed
}

func NewOrientationView() *OrientationView {
	return &OrientationView{}
}

func (v *OrientationView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.String = "Get orientation"
	chl1.OnPress = func() {
		o := application.Orientation()
		view.Alert("Orientation"+strconv.Itoa(int(o)), "")
	}
	_ = l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})
	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

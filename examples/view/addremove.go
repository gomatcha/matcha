package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewAddRemoveView", func() *view.Root {
		return view.NewRoot(NewAddRemoveView())
	})
}

type AddRemoveView struct {
	view.Embed
	showView bool
}

func NewAddRemoveView() *AddRemoveView {
	return &AddRemoveView{}
}

func (v *AddRemoveView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.Enabled = true
	chl1.String = "Press Me"
	chl1.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	chl1.OnPress = func() {
		v.showView = !v.showView
		v.Signal()
	}
	_ = l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(0)
		s.Width(200)
	})

	if v.showView {
		chl2 := view.NewBasicView()
		chl2.Painter = &paint.Style{BackgroundColor: colornames.Green}
		l.Add(chl2, func(s *constraint.Solver) {
			s.Top(200)
			s.Left(0)
			s.Width(100)
			s.Height(100)
		})
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

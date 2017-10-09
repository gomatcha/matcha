package view

import (
	"runtime"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewAddRemoveView", func() view.View {
		return NewAddRemoveView()
	})
}

type AddRemoveView struct {
	view.Embed
	showView bool
}

func NewAddRemoveView() *AddRemoveView {
	return &AddRemoveView{}
}

func (v *AddRemoveView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.Enabled = true
	if !v.showView {
		chl1.String = "Add"
	} else {
		chl1.String = "Remove"
	}
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

		var chl3 view.View
		if runtime.GOOS == "android" {
			chl3 = android.NewStackView()
		} else {
			chl3 = ios.NewStackView()
		}
		l.Add(chl3, func(s *constraint.Solver) {
			s.Top(200)
			s.Left(100)
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

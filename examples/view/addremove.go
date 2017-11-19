package view

import (
	"runtime"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
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

	label := view.NewTextView()
	label.String = "Toggle view:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	button := view.NewButton()
	if !v.showView {
		button.String = "Add"
	} else {
		button.String = "Remove"
	}
	button.OnPress = func() {
		v.showView = !v.showView
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	if v.showView {
		chl1 := view.NewBasicView()
		chl1.Painter = &paint.Style{BackgroundColor: colornames.Green}
		g = l.Add(chl1, func(s *constraint.Solver) {
			s.TopEqual(g.Bottom())
			s.LeftEqual(g.Left())
			s.Width(100)
			s.Height(100)
		})

		var chl2 view.View
		if runtime.GOOS == "android" {
			chl2 = android.NewPagerView()
		} else {
			chl2 = ios.NewStackView()
		}
		g = l.Add(chl2, func(s *constraint.Solver) {
			s.TopEqual(g.Bottom())
			s.LeftEqual(g.Left())
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

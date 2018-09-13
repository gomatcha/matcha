package android

import (
	"image/color"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"github.com/gomatcha/matcha/view/android"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view/android NewStatusBarView", func() view.View {
		return NewStatusBarView()
	})
}

type StatusBarView struct {
	view.Embed
	color color.Color
	style android.StatusBarStyle
}

func NewStatusBarView() *StatusBarView {
	return &StatusBarView{}
}

func (v *StatusBarView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewButton()
	chl1.String = "Toggle Color"
	chl1.OnPress = func() {
		if v.color == colornames.Red {
			v.color = colornames.White
		} else {
			v.color = colornames.Red
		}
		v.Signal()
	}
	l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := view.NewButton()
	chl2.String = "Toggle Style"
	chl2.OnPress = func() {
		if v.style == android.StatusBarStyleLight {
			v.style = android.StatusBarStyleDark
		} else {
			v.style = android.StatusBarStyleLight
		}
		v.Signal()
	}
	l.Add(chl2, func(s *constraint.Solver) {
		s.Top(200)
		s.Left(100)
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			&android.StatusBar{
				Color: v.color,
				Style: v.style,
			},
		},
	}
}

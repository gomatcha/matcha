package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/android NewSwipeView", func() *view.Root {
		return view.NewRoot(NewSwipeView())
	})
}

type SwipeView struct {
	view.Embed
}

func NewSwipeView() *SwipeView {
	return &SwipeView{}
}

func (v *SwipeView) Build(ctx *view.Context) view.Model {
	chl := android.NewSwipeView()

	return view.Model{
		Children: []view.View{chl},
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

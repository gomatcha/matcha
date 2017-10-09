package android

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/android NewPagerView", func() view.View {
		return NewPagerView()
	})
}

func NewPagerView() view.View {
	v := android.NewPagerView()

	app := &PagerApp{}
	app.Pages = v.Pages
	app.Pages.Notify(func() {
		fmt.Println("CurrentPage", v.Pages.SelectedIndex())
	})

	v1 := NewPagerChildView()
	v1.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	v1.PagerButton = &android.PagerButton{Title: "Title 1"}

	v2 := NewPagerChildView()
	v2.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	v2.PagerButton = &android.PagerButton{Title: "Title 2"}

	v3 := NewPagerChildView()
	v3.PaintStyle = &paint.Style{BackgroundColor: colornames.Black}
	v3.PagerButton = &android.PagerButton{Title: "Title 3"}

	app.Pages.SetViews(v1, v2, v3)
	app.Pages.SetSelectedIndex(2)

	return v
}

type PagerApp struct {
	Pages *android.Pages
}

type PagerChildView struct {
	view.Embed
	PaintStyle  *paint.Style
	PagerButton *android.PagerButton
}

func NewPagerChildView() *PagerChildView {
	return &PagerChildView{}
}

func (v *PagerChildView) Build(ctx view.Context) view.Model {
	var p paint.Painter
	if v.PaintStyle != nil {
		p = v.PaintStyle
	}
	return view.Model{
		Painter: p,
		Options: []view.Option{
			v.PagerButton,
		},
	}
}

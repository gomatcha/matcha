// Package tabview provides examples of how to use the matcha/view/tabview package.
package tabview

import (
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/tabview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/tabview New", func() *view.Root {
		app := &App{tabs: &tabview.Tabs{}}

		view1 := NewTouchView(app)
		view1.Color = colornames.Blue
		button1 := &tabview.Button{
			Title: "Title 1",
			Badge: "badge",
			// Icon:         env.MustLoadImage("TabCamera"),
			// SelectedIcon: env.MustLoadImage("TabCameraFilled"),
		}

		view2 := NewTouchView(app)
		view2.Color = colornames.Red
		button2 := &tabview.Button{
			Title: "Title 2",
			// Icon:         env.MustLoadImage("TabMap"),
			// SelectedIcon: env.MustLoadImage("TabMapFilled"),
		}

		view3 := NewTouchView(app)
		view3.Color = colornames.Yellow

		view4 := NewTouchView(app)
		view4.Color = colornames.Green

		v := tabview.New()
		v.BarColor = colornames.White
		v.SelectedColor = colornames.Red
		v.UnselectedColor = colornames.Darkgray
		v.Tabs = app.tabs
		v.Tabs.SetSelectedIndex(1)
		v.Tabs.SetViews(
			tabview.WithButton(view1, button1),
			tabview.WithButton(view2, button2),
			view3,
			view4,
		)
		return view.NewRoot(v)
	})
}

type App struct {
	tabs *tabview.Tabs
}

type TouchView struct {
	view.Embed
	app    *App
	Color  color.Color
	button *tabview.Button
}

func NewTouchView(app *App) *TouchView {
	return &TouchView{
		app: app,
		button: &tabview.Button{
			Title: "Testing",
			// Icon:         env.MustLoadImage("TabSearch"),
			// SelectedIcon: env.MustLoadImage("TabSearchFilled"),
		},
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapRecognizer{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			v.app.tabs.SetSelectedIndex(0)
			v.button.Title = "Updated"
			v.Signal()
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
	}
}

func (v *TouchView) TabButton(*view.Context) *tabview.Button {
	return v.button
}

package ios

import (
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewTabView", func() *view.Root {
		app := &App{tabs: &ios.Tabs{}}

		view1 := NewTouchView(app)
		view1.Color = colornames.Blue
		view1.button = &ios.TabButton{
			Title: "Title 1",
			Badge: "badge",
			// Icon:         env.MustLoadImage("TabCamera"),
			// SelectedIcon: env.MustLoadImage("TabCameraFilled"),
		}

		view2 := NewTouchView(app)
		view2.Color = colornames.Red
		view2.button = &ios.TabButton{
			Title: "Title 2",
			// Icon:         env.MustLoadImage("TabMap"),
			// SelectedIcon: env.MustLoadImage("TabMapFilled"),
		}

		view3 := NewTouchView(app)
		view3.Color = colornames.Yellow

		view4 := NewTouchView(app)
		view4.Color = colornames.Green

		v := ios.NewTabView()
		v.BarColor = colornames.White
		v.SelectedColor = colornames.Red
		v.UnselectedColor = colornames.Darkgray
		v.Tabs = app.tabs
		v.Tabs.SetSelectedIndex(1)
		v.Tabs.SetViews(
			view1,
			view2,
			view3,
			view4,
		)
		return view.NewRoot(v)
	})
}

type App struct {
	tabs *ios.Tabs
}

type TouchView struct {
	view.Embed
	app    *App
	Color  color.Color
	button *ios.TabButton
}

func NewTouchView(app *App) *TouchView {
	return &TouchView{
		app: app,
		button: &ios.TabButton{
			Title: "Testing",
			// Icon:         env.MustLoadImage("TabSearch"),
			// SelectedIcon: env.MustLoadImage("TabSearchFilled"),
		},
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapGesture{
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
			touch.GestureList{tap},
			v.button,
		},
	}
}

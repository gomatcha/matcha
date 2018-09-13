package ios

import (
	"image/color"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/pointer"
	"github.com/gomatcha/matcha/view"
	"github.com/gomatcha/matcha/view/ios"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view/ios NewTabView", func() view.View {
		return NewTabView()
	})
}

func NewTabView() view.View {
	app := &TabApp{tabs: &ios.Tabs{}}

	view1 := NewTabChild(app)
	view1.Color = colornames.Blue
	view1.button = &ios.TabButton{
		Title: "Title 1",
		Badge: "badge",
		// Icon:         env.MustLoadImage("TabCamera"),
		// SelectedIcon: env.MustLoadImage("TabCameraFilled"),
	}

	view2 := NewTabChild(app)
	view2.Color = colornames.Red
	view2.button = &ios.TabButton{
		Title: "Title 2",
		// Icon:         env.MustLoadImage("TabMap"),
		// SelectedIcon: env.MustLoadImage("TabMapFilled"),
	}

	view3 := NewTabChild(app)
	view3.Color = colornames.Yellow

	view4 := NewTabChild(app)
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
	return v
}

type TabApp struct {
	tabs *ios.Tabs
}

type TabChild struct {
	view.Embed
	app    *TabApp
	Color  color.Color
	button *ios.TabButton
}

func NewTabChild(app *TabApp) *TabChild {
	return &TabChild{
		app: app,
		button: &ios.TabButton{
			Title: "Testing",
			// Icon:         env.MustLoadImage("TabSearch"),
			// SelectedIcon: env.MustLoadImage("TabSearchFilled"),
		},
	}
}

func (v *TabChild) Build(ctx view.Context) view.Model {
	tap := &pointer.TapGesture{
		Count: 1,
		OnEvent: func(e *pointer.TapEvent) {
			v.app.tabs.SetSelectedIndex(0)
			v.button.Title = "Updated"
			v.Signal()
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			pointer.GestureList{tap},
			v.button,
		},
	}
}

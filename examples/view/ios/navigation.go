package ios

import (
	"image/color"

	"golang.org/x/image/colornames"

	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewNavigationView", func() view.View {
		return NewNavigationView()
	})
}

func NewNavigationView() view.View {
	stackview1 := ios.NewStackView()
	stackview2 := ios.NewStackView()
	stackview3 := ios.NewStackView()
	stackview4 := ios.NewStackView()
	tabview := ios.NewTabView()

	app := &NavigationApp{
		stack1: stackview1.Stack,
		stack2: stackview2.Stack,
		stack3: stackview3.Stack,
		stack4: stackview4.Stack,
		tabs:   tabview.Tabs,
	}
	app.stack1.SetViews(NewNavigationChild(app))
	app.stack2.SetViews(NewNavigationChild(app))
	app.stack3.SetViews(NewNavigationChild(app))
	app.stack4.SetViews(NewNavigationChild(app))
	app.tabs.SetViews(stackview1, stackview2, stackview3, stackview4)
	return tabview
}

type NavigationApp struct {
	tabs   *ios.Tabs
	stack1 *ios.Stack
	stack2 *ios.Stack
	stack3 *ios.Stack
	stack4 *ios.Stack
}

func (app *NavigationApp) CurrentStackView() *ios.Stack {
	switch app.tabs.SelectedIndex() {
	case 0:
		return app.stack1
	case 1:
		return app.stack2
	case 2:
		return app.stack3
	case 3:
		return app.stack4
	}
	return nil
}

type NavigationChild struct {
	view.Embed
	app   *NavigationApp
	Color color.Color
}

func NewNavigationChild(app *NavigationApp) *NavigationChild {
	return &NavigationChild{
		Color: colornames.White,
		app:   app,
	}
}

func (v *NavigationChild) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			&pointer.TapGesture{
				Count: 1,
				OnRecognize: func(e *pointer.TapEvent) {
					child := NewNavigationChild(v.app)
					child.Color = colornames.Red
					v.app.CurrentStackView().Push(child)
				},
			},
		},
	}
}

// Package screen provides examples of how to add navigation components.
package screen

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/screen New", func() *view.Root {
		stackview1 := ios.NewStackView()
		stackview2 := ios.NewStackView()
		stackview3 := ios.NewStackView()
		stackview4 := ios.NewStackView()
		tabview := ios.NewTabView()

		app := &App{
			stack1: stackview1.Stack,
			stack2: stackview2.Stack,
			stack3: stackview3.Stack,
			stack4: stackview4.Stack,
			tabs:   tabview.Tabs,
		}
		app.stack1.SetViews(NewTouchView(app))
		app.stack2.SetViews(NewTouchView(app))
		app.stack3.SetViews(NewTouchView(app))
		app.stack4.SetViews(NewTouchView(app))
		app.tabs.SetViews(stackview1, stackview2, stackview3, stackview4)
		return view.NewRoot(tabview)
	})
}

type App struct {
	tabs   *ios.Tabs
	stack1 *ios.Stack
	stack2 *ios.Stack
	stack3 *ios.Stack
	stack4 *ios.Stack
}

func (app *App) CurrentStackView() *ios.Stack {
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

type TouchView struct {
	view.Embed
	app   *App
	Color color.Color
}

func NewTouchView(app *App) *TouchView {
	return &TouchView{
		Color: colornames.White,
		app:   app,
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapGesture{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			child := NewTouchView(v.app)
			child.Color = colornames.Red
			v.app.CurrentStackView().Push(child)
			fmt.Println("child", child)
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			touch.GestureList{tap},
		},
	}
}

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
	"gomatcha.io/matcha/view/ios/stackview"
	"gomatcha.io/matcha/view/ios/tabview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/screen New", func() *view.Root {
		app := &App{
			stack1: &stackview.Stack{},
			stack2: &stackview.Stack{},
			stack3: &stackview.Stack{},
			stack4: &stackview.Stack{},
			tabs:   &tabview.Tabs{},
		}

		// Configure the stacks
		app.stack1.SetViews(NewTouchView(nil, "", app))
		app.stack2.SetViews(NewTouchView(nil, "", app))
		app.stack3.SetViews(NewTouchView(nil, "", app))
		app.stack4.SetViews(NewTouchView(nil, "", app))

		// Configure the tabs
		stackview1 := stackview.New()
		stackview1.Stack = app.stack1
		stackview2 := stackview.New()
		stackview2.Stack = app.stack2
		stackview3 := stackview.New()
		stackview3.Stack = app.stack3
		stackview4 := stackview.New()
		stackview4.Stack = app.stack4
		app.tabs.SetViews(
			stackview1,
			stackview2,
			stackview3,
			stackview4,
		)

		// Return tabview
		v := tabview.New()
		v.Tabs = app.tabs
		return view.NewRoot(v)
	})
}

type App struct {
	tabs   *tabview.Tabs
	stack1 *stackview.Stack
	stack2 *stackview.Stack
	stack3 *stackview.Stack
	stack4 *stackview.Stack
}

func (app *App) CurrentStackView() *stackview.Stack {
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

func NewTouchView(ctx *view.Context, key string, app *App) *TouchView {
	return &TouchView{
		Color: colornames.White,
		app:   app,
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapRecognizer{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			child := NewTouchView(nil, "", v.app)
			child.Color = colornames.Red
			v.app.CurrentStackView().Push(child)
			fmt.Println("child", child)
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
	}
}

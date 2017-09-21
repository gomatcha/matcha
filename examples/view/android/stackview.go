package android

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/android NewStackView", func() view.View {
		stackview := android.NewStackView()
		app := &StackApp{
			stack: stackview.Stack,
		}

		view1 := NewStackChild(app)
		view1.Color = colornames.White
		v1 := view.WithOptions(view1, &android.StackBar{Title: "Title 1"})

		// view2 := NewStackChild(app)
		// view2.Color = colornames.Red
		// v2 := view.WithOptions(view2, &android.StackBar{Title: "Title 2"})

		// view3 := NewStackChild(app)
		// view3.Color = colornames.Yellow

		// view4 := NewStackChild(app)
		// view4.Color = colornames.Green

		// app.stack.SetViews(v1, v2, view3, view4)
		app.stack.SetViews(v1)
		return stackview
	})
}

type StackApp struct {
	stack *android.Stack
}

type StackChild struct {
	view.Embed
	app   *StackApp
	Color color.Color
	Index int
}

func NewStackChild(app *StackApp) *StackChild {
	return &StackChild{
		app: app,
	}
}

func (v *StackChild) Build(ctx view.Context) view.Model {
	// tap := &touch.NewTapGesture()
	// tap.OnTouch = func(e *touch.TapEvent) {
	// 	if e.Kind != touch.EventKindRecognized {
	// 		return
	// 	}
	// 	// v.bar.Title = "Updated"
	// 	// v.Signal()

	// 	child := NewStackChild(v.app)
	// 	child.Index = v.Index + 1
	// 	child.Color = colornames.White
	// 	v.app.stack.Push(child)
	// }

	tap := &pointer.TapGesture{
		Count: 1,
		OnEvent: func(e *pointer.TapEvent) {
			if e.Kind != pointer.EventKindRecognized {
				return
			}
			// v.bar.Title = "Updated"
			// v.Signal()

			child := NewStackChild(v.app)
			child.Index = v.Index + 1
			child.Color = colornames.White
			v.app.stack.Push(child)
		},
	}

	titleStyle := &text.Style{}
	title := text.NewStyledText(fmt.Sprintf("Title %v", v.Index), titleStyle)
	subtitleStyle := &text.Style{}
	subtitle := text.NewStyledText("Subtitle", subtitleStyle)

	item := android.NewStackBarItem()
	item.Title = "index"
	item.Icon = application.MustLoadImage("settings_airplane")
	item.OnPress = func() {
		fmt.Println("OnPress")
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			pointer.GestureList{tap},
			&android.StackBar{
				StyledTitle:    title,
				StyledSubtitle: subtitle,
				Color:          colornames.White,
				Items:          []*android.StackBarItem{item},
			},
		},
	}
}

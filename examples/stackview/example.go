// Package stackview provides examples of how to use the matcha/view/stackview package.
package stackview

import (
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
	"gomatcha.io/matcha/view/stackview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/stackview New", func() *view.Root {
		app := &App{
			stack: &stackview.Stack{},
		}

		view1 := NewTouchView(nil, "", app)
		view1.Color = colornames.Blue
		bar1 := &stackview.Bar{
			Title: "Title 1",
		}

		// view2 := NewTouchView(nil, "", app)
		// view2.Color = colornames.Red
		// bar2 := &stackview.Bar{
		// 	Title: "Title 2",
		// }

		// view3 := NewTouchView(nil, "", app)
		// view3.Color = colornames.Yellow
		// view4 := NewTouchView(nil, "", app)
		// view4.Color = colornames.Green

		v := stackview.New()
		v.Stack = app.stack
		v.Stack.SetViews(
			stackview.WithBar(view1, bar1),
			// stackview.WithBar(view2, bar2),
			// view3,
			// view4,
		)
		return view.NewRoot(v)
	})
}

type App struct {
	stack *stackview.Stack
}

type TouchView struct {
	view.Embed
	app   *App
	Color color.Color
	bar   *stackview.Bar
}

func NewTouchView(ctx *view.Context, key string, app *App) *TouchView {
	if v, ok := ctx.Prev(key).(*TouchView); ok {
		return v
	}
	return &TouchView{
		Embed: view.Embed{Key: key},
		app:   app,
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapRecognizer{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			// v.bar.Title = "Updated"
			// v.Signal()

			child := NewTouchView(nil, "", v.app)
			child.Color = colornames.Purple
			v.app.stack.Push(child)
		},
	}

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
	}
}

func (v *TouchView) StackBar(ctx *view.Context) *stackview.Bar {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(100))
		s.WidthEqual(constraint.Const(100))
	})

	titleView := basicview.New()
	titleView.Painter = &paint.Style{BackgroundColor: colornames.Red}
	titleView.Layouter = l

	l2 := &constraint.Layouter{}
	l2.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	rightView := basicview.New()
	rightView.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	rightView.Layouter = l2

	l3 := &constraint.Layouter{}
	l3.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	leftView := basicview.New()
	leftView.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	leftView.Layouter = l3

	return &stackview.Bar{
		Title:      "Title",
		TitleView:  titleView,
		RightViews: []view.View{rightView},
		LeftViews:  []view.View{leftView},
	}
}

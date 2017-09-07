package ios

import (
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewStackView", func() *view.Root {
		stackview := ios.NewStackView()
		app := &StackApp{
			stack: stackview.Stack,
		}

		view1 := NewStackChild(app)
		view1.Color = colornames.Blue
		v1 := view.WithOptions(view1, &ios.StackBar{Title: "Title 1"})

		view2 := NewStackChild(app)
		view2.Color = colornames.Red
		v2 := view.WithOptions(view2, &ios.StackBar{Title: "Title 2"})

		view3 := NewStackChild(app)
		view3.Color = colornames.Yellow

		view4 := NewStackChild(app)
		view4.Color = colornames.Green

		app.stack.SetViews(v1, v2, view3, view4)
		return view.NewRoot(stackview)
	})
}

type StackApp struct {
	stack *ios.Stack
}

type StackChild struct {
	view.Embed
	app   *StackApp
	Color color.Color
	bar   *ios.StackBar
}

func NewStackChild(app *StackApp) *StackChild {
	return &StackChild{
		app: app,
	}
}

func (v *StackChild) Build(ctx *view.Context) view.Model {
	tap := &touch.TapGesture{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			// v.bar.Title = "Updated"
			// v.Signal()

			child := NewStackChild(v.app)
			child.Color = colornames.Purple
			v.app.stack.Push(child)
		},
	}

	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(100))
		s.WidthEqual(constraint.Const(100))
	})

	titleView := view.NewBasicView()
	titleView.Painter = &paint.Style{BackgroundColor: colornames.Red}
	titleView.Layouter = l

	l2 := &constraint.Layouter{}
	l2.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	rightView := view.NewBasicView()
	rightView.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	rightView.Layouter = l2

	l3 := &constraint.Layouter{}
	l3.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	leftView := view.NewBasicView()
	leftView.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	leftView.Layouter = l3

	return view.Model{
		Painter: &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			touch.GestureList{tap},
			&ios.StackBar{
				Title:      "Title",
				TitleView:  titleView,
				RightViews: []view.View{rightView},
				LeftViews:  []view.View{leftView},
			},
		},
	}
}

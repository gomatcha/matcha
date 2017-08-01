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
	bridge.RegisterFunc("gomatcha.io/matcha/examples/stackscreen New", func() *view.Root {
		return view.NewRoot(NewAppView())
	})
}

type App struct {
	stack *stackview.Stack
}

func NewAppView() view.View {
	app := &App{}

	screen1 := NewTouchScreen(app, colornames.Blue)
	bar1 := &stackview.Bar{
		Title: "Title 1",
	}

	screen2 := NewTouchScreen(app, colornames.Red)
	bar2 := &stackview.Bar{
		Title: "Title 2",
	}

	screen3 := NewTouchScreen(app, colornames.Yellow)
	screen4 := NewTouchScreen(app, colornames.Green)

	app.stack = &stackview.Stack{}
	app.stack.SetViews(
		stackview.WithBar(screen1, bar1),
		stackview.WithBar(screen2, bar2),
		screen3,
		screen4,
	)

	v := stackview.New(nil, "")
	v.Stack = app.stack
	return v
}

func NewTouchScreen(app *App, c color.Color) view.View {
	chl := NewTouchView(nil, "", app)
	chl.Color = c
	return chl
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
		Embed: ctx.NewEmbed(key),
		app:   app,
	}
}

func (v *TouchView) Build(ctx *view.Context) view.Model {
	tap := &touch.TapRecognizer{
		Count: 1,
		OnTouch: func(e *touch.TapEvent) {
			// v.bar.Title = "Updated"
			// v.Signal()

			v.app.stack.Push(NewTouchScreen(v.app, colornames.Purple))
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

	titleView := basicview.New(ctx, "axbaba")
	titleView.Painter = &paint.Style{BackgroundColor: colornames.Red}
	titleView.Layouter = l

	l2 := &constraint.Layouter{}
	l2.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	rightView := basicview.New(ctx, "right")
	rightView.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	rightView.Layouter = l2

	l3 := &constraint.Layouter{}
	l3.Solve(func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.HeightEqual(constraint.Const(50))
		s.WidthEqual(constraint.Const(50))
	})
	leftView := basicview.New(ctx, "left")
	leftView.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	leftView.Layouter = l3

	return &stackview.Bar{
		Title:      "Title",
		TitleView:  titleView,
		RightViews: []view.View{rightView},
		LeftViews:  []view.View{leftView},
	}
}

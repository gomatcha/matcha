package todo

import (
	"image/color"
	"strconv"

	"golang.org/x/image/colornames"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
	"gomatcha.io/matcha/view/imageview"
	"gomatcha.io/matcha/view/scrollview"
	"gomatcha.io/matcha/view/stackview"
	"gomatcha.io/matcha/view/textinput"
	"gomatcha.io/matcha/view/textview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/todo New", func() *view.Root {
		app := NewAppView(nil, "")

		stack := &stackview.Stack{}
		stack.SetViews(app)

		titleTextStyle := &text.Style{}
		titleTextStyle.SetFont(text.Font{
			Family: "Helvetica Neue",
			Face:   "Medium",
			Size:   20,
		})
		titleTextStyle.SetTextColor(colornames.White)

		v := stackview.New(nil, "")
		v.Stack = stack
		v.BarColor = color.RGBA{R: 46, G: 124, B: 190, A: 1}
		v.TitleTextStyle = titleTextStyle
		return view.NewRoot(v)
	})
}

type Todo struct {
	Title     string
	Completed bool
}

type AppView struct {
	view.Embed
	Todos []*Todo
}

func NewAppView(ctx *view.Context, key string) *AppView {
	if v, ok := ctx.Prev(key).(*AppView); ok {
		return v
	}
	return &AppView{Embed: ctx.NewEmbed(key)}
}

func (v *AppView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}

	for i, todo := range v.Todos {
		idx := i
		todoView := NewTodoView(ctx, strconv.Itoa(idx))
		todoView.Todo = todo
		todoView.OnDelete = func() {
			v.Todos = append(v.Todos[:idx], v.Todos[idx+1:]...)
			v.Signal()
		}
		todoView.OnComplete = func(complete bool) {
			v.Todos[idx].Completed = complete
			v.Signal()
		}
		l.Add(todoView, nil)
	}

	addView := NewAddView(ctx, "add")
	addView.OnAdd = func(title string) {
		v.Todos = append(v.Todos, &Todo{Title: title})
		v.Signal()
	}
	l.Add(addView, nil)

	scrollView := scrollview.New(ctx, "scrollView")
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l
	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			app.StatusBar{Style: app.StatusBarStyleLight},
		},
	}
}

func (v *AppView) StackBar(ctx *view.Context) *stackview.Bar {
	return &stackview.Bar{Title: "To Do Example"}
}

type AddView struct {
	view.Embed
	text      *text.Text
	responder keyboard.Responder
	OnAdd     func(title string)
}

func NewAddView(ctx *view.Context, key string) *AddView {
	if v, ok := ctx.Prev(key).(*AddView); ok {
		return v
	}
	return &AddView{
		Embed: ctx.NewEmbed(key),
		text:  text.New(""),
	}
}

func (v *AddView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
		s.WidthEqual(l.MaxGuide().Width())
	})

	style := &text.Style{}
	style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   20,
	})

	placeholderStyle := &text.Style{}
	placeholderStyle.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   20,
	})
	placeholderStyle.SetTextColor(colornames.Lightgray)

	input := textinput.New(ctx, "input")
	input.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	input.Text = v.text
	input.Style = style
	input.PlaceholderText = text.New("What needs to be done?")
	input.PlaceholderStyle = placeholderStyle
	input.KeyboardReturnType = keyboard.DoneReturnType
	input.Responder = &v.responder
	input.OnSubmit = func() {
		str := v.text.String()
		v.responder.Dismiss()
		v.text.SetString("")
		if str != "" {
			v.OnAdd(str)
		}
	}
	l.Add(input, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(15))
		s.RightEqual(l.Right().Add(-15))
		s.CenterYEqual(l.CenterY())
	})

	separator := basicview.New(ctx, "separator")
	separator.Painter = &paint.Style{BackgroundColor: color.RGBA{203, 202, 207, 255}}
	l.Add(separator, func(s *constraint.Solver) {
		s.Height(1)
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.BottomEqual(l.Bottom())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
	}
}

type TodoView struct {
	view.Embed
	Todo       *Todo
	OnDelete   func()
	OnComplete func(check bool)
}

func NewTodoView(ctx *view.Context, key string) *TodoView {
	if v, ok := ctx.Prev(key).(*TodoView); ok {
		return v
	}
	return &TodoView{Embed: ctx.NewEmbed(key)}
}

func (v *TodoView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
		s.WidthEqual(l.MaxGuide().Width())
	})

	checkbox := NewCheckbox(ctx, "checkbox")
	checkbox.Value = v.Todo.Completed
	checkbox.OnValueChange = func(value bool) {
		v.OnComplete(value)
	}
	checkboxGuide := l.Add(checkbox, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(l.Left().Add(15))
	})

	deleteButton := NewDeleteButton(ctx, "delete")
	deleteButton.OnPress = func() {
		v.OnDelete()
	}
	deleteGuide := l.Add(deleteButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.RightEqual(l.Right().Add(-15))
	})

	titleView := textview.New(ctx, "title")
	titleView.String = v.Todo.Title
	titleView.Style = nil //...
	l.Add(titleView, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(checkboxGuide.Right().Add(15))
		s.RightEqual(deleteGuide.Left().Add(-15))
	})

	separator := basicview.New(ctx, "separator")
	separator.Painter = &paint.Style{BackgroundColor: color.RGBA{203, 202, 207, 255}}
	l.Add(separator, func(s *constraint.Solver) {
		s.Height(1)
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.BottomEqual(l.Bottom())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
	}
}

type Checkbox struct {
	view.Embed
	Value         bool
	OnValueChange func(value bool)
}

func NewCheckbox(ctx *view.Context, key string) *Checkbox {
	if v, ok := ctx.Prev(key).(*Checkbox); ok {
		return v
	}
	return &Checkbox{Embed: ctx.NewEmbed(key)}
}

func (v *Checkbox) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(40)
		s.Height(40)
	})

	imageView := imageview.New(ctx, "image")
	if v.Value {
		imageView.Image = app.MustLoadImage("CheckboxChecked")
	} else {
		imageView.Image = app.MustLoadImage("CheckboxUnchecked")
	}
	l.Add(imageView, func(s *constraint.Solver) {
		s.CenterXEqual(l.CenterX())
		s.CenterYEqual(l.CenterY())
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
	})

	button := &touch.ButtonRecognizer{
		OnTouch: func(e *touch.ButtonEvent) {
			if e.Kind == touch.EventKindRecognized {
				v.OnValueChange(!v.Value)
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		// Painter:  painter,
		Layouter: l,
		Options: []view.Option{
			touch.RecognizerList{button},
		},
	}
}

type DeleteButton struct {
	view.Embed
	OnPress func()
}

func NewDeleteButton(ctx *view.Context, key string) *DeleteButton {
	if v, ok := ctx.Prev(key).(*DeleteButton); ok {
		return v
	}
	return &DeleteButton{Embed: ctx.NewEmbed(key)}
}

func (v *DeleteButton) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(40)
		s.Height(40)
	})

	imageView := imageview.New(ctx, "image")
	imageView.Image = app.MustLoadImage("Delete")
	l.Add(imageView, func(s *constraint.Solver) {
		s.CenterXEqual(l.CenterX())
		s.CenterYEqual(l.CenterY())
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
	})

	button := &touch.ButtonRecognizer{
		OnTouch: func(e *touch.ButtonEvent) {
			if e.Kind == touch.EventKindRecognized {
				v.OnPress()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			touch.RecognizerList{button},
		},
	}
}

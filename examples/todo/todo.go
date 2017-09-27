// Package todo provides an example of a basic Todo app.
package todo

import (
	"image/color"
	"runtime"

	"golang.org/x/image/colornames"

	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/todo New", func() view.View {
		appview := NewAppView()
		appview.Todos = []*Todo{
			&Todo{Title: "Eat Chicken"},
			&Todo{Title: "Kill Dog"},
		}

		if runtime.GOOS == "ios" {
			v := ios.NewStackView()
			v.Stack.SetViews(appview)
			v.BarColor = color.RGBA{R: 46, G: 124, B: 190, A: 1}
			v.TitleTextStyle = &text.Style{}
			v.TitleTextStyle.SetFont(text.FontWithName("HelveticaNeue-Medium", 20))
			v.TitleTextStyle.SetTextColor(colornames.White)
			return v
		} else {
			v := android.NewStackView()
			v.Stack.SetViews(appview)
			// v.BarColor = color.RGBA{R: 46, G: 124, B: 190, A: 1}
			// v.TitleTextStyle = &text.Style{}
			// v.TitleTextStyle.SetFont(text.FontWithName("HelveticaNeue-Medium", 20))
			// v.TitleTextStyle.SetTextColor(colornames.White)
			return v
		}
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

func NewAppView() *AppView {
	return &AppView{}
}

func (v *AppView) Build(ctx view.Context) view.Model {
	l := &table.Layouter{}

	for i, todo := range v.Todos {
		idx := i
		todoView := NewTodoView()
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

	addView := NewAddView()
	addView.OnAdd = func(title string) {
		v.Todos = append(v.Todos, &Todo{Title: title})
		v.Signal()
	}
	l.Add(addView, nil)

	scrollView := view.NewScrollView()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l
	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			&ios.StackBar{Title: "Todos"},
			&ios.StatusBar{Style: ios.StatusBarStyleLight},
			&android.StackBar{Title: "Todos"},
		},
	}
}

type AddView struct {
	view.Embed
	text      *text.Text
	responder keyboard.Responder
	OnAdd     func(title string)
}

func NewAddView() *AddView {
	return &AddView{
		text: text.New(""),
	}
}

func (v *AddView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
		s.WidthEqual(l.MaxGuide().Width())
	})

	style := &text.Style{}
	style.SetFont(text.FontWithName("HelveticaNeue", 20))

	placeholderStyle := &text.Style{}
	placeholderStyle.SetFont(text.FontWithName("HelveticaNeue", 20))
	placeholderStyle.SetTextColor(colornames.Lightgray)

	input := view.NewTextInput()
	input.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
	input.Text = v.text
	input.Style = style
	input.Placeholder = "What needs to be done?"
	input.PlaceholderStyle = placeholderStyle
	input.Responder = &v.responder
	input.OnSubmit = func(t *text.Text) {
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

	separator := view.NewBasicView()
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
		// Painter:  &paint.Style{BackgroundColor: colornames.Red},
	}
}

type TodoView struct {
	view.Embed
	Todo       *Todo
	OnDelete   func()
	OnComplete func(check bool)
}

func NewTodoView() *TodoView {
	return &TodoView{}
}

func (v *TodoView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
		s.WidthEqual(l.MaxGuide().Width())
	})

	checkbox := NewCheckbox()
	checkbox.Value = v.Todo.Completed
	checkbox.OnValueChange = func(value bool) {
		v.OnComplete(value)
	}
	checkboxGuide := l.Add(checkbox, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(l.Left().Add(15))
	})

	deleteButton := NewDeleteButton()
	deleteButton.OnPress = func() {
		v.OnDelete()
	}
	deleteGuide := l.Add(deleteButton, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.RightEqual(l.Right().Add(-15))
	})

	titleView := view.NewTextView()
	titleView.String = v.Todo.Title
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	l.Add(titleView, func(s *constraint.Solver) {
		s.CenterYEqual(l.CenterY())
		s.LeftEqual(checkboxGuide.Right().Add(15))
		s.RightEqual(deleteGuide.Left().Add(-15))
	})

	separator := view.NewBasicView()
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

func NewCheckbox() *Checkbox {
	return &Checkbox{}
}

func (v *Checkbox) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(40)
		s.Height(40)
	})

	imageView := view.NewImageView()
	if v.Value {
		imageView.Image = application.MustLoadImage("checkbox_checked")
	} else {
		imageView.Image = application.MustLoadImage("checkbox_unchecked")
	}
	l.Add(imageView, func(s *constraint.Solver) {
		s.CenterXEqual(l.CenterX())
		s.CenterYEqual(l.CenterY())
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
	})

	button := &pointer.ButtonGesture{
		OnEvent: func(e *pointer.ButtonEvent) {
			if e.Kind == pointer.EventKindRecognized {
				v.OnValueChange(!v.Value)
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		// Painter:  painter,
		Layouter: l,
		Options: []view.Option{
			pointer.GestureList{button},
		},
	}
}

type DeleteButton struct {
	view.Embed
	OnPress func()
}

func NewDeleteButton() *DeleteButton {
	return &DeleteButton{}
}

func (v *DeleteButton) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(40)
		s.Height(40)
	})

	imageView := view.NewImageView()
	imageView.Image = application.MustLoadImage("delete")
	l.Add(imageView, func(s *constraint.Solver) {
		s.CenterXEqual(l.CenterX())
		s.CenterYEqual(l.CenterY())
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
	})

	button := &pointer.ButtonGesture{
		OnEvent: func(e *pointer.ButtonEvent) {
			if e.Kind == pointer.EventKindRecognized {
				v.OnPress()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			pointer.GestureList{button},
		},
	}
}

package ios

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewStackView", func() view.View {
		return NewStackView()
	})
}

type StackApp struct {
	Stack             *ios.Stack
	BarColor          comm.ColorValue
	ItemTintColor     comm.ColorValue
	ItemTitleStyle    comm.InterfaceValue
	TitleStyle        comm.InterfaceValue
	AllItemTintColor  comm.ColorValue
	AllItemTitleStyle comm.InterfaceValue
}

func NewStackView() view.View {
	app := &StackApp{
		Stack: &ios.Stack{},
	}

	view1 := NewStackConfigureView()
	view1.App = app
	app.Stack.SetViews(view1)
	return &StackAppView{
		App: app,
	}
}

type StackAppView struct {
	view.Embed
	App *StackApp
}

func (v *StackAppView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(&v.App.BarColor)
		v.Subscribe(&v.App.TitleStyle)
		v.Subscribe(&v.App.AllItemTintColor)
		v.Subscribe(&v.App.AllItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.App.BarColor)
		v.Unsubscribe(&v.App.TitleStyle)
		v.Unsubscribe(&v.App.AllItemTintColor)
		v.Unsubscribe(&v.App.AllItemTitleStyle)
	}
}

func (v *StackAppView) Build(ctx view.Context) view.Model {
	stackview := ios.NewStackView()
	stackview.Stack = v.App.Stack
	stackview.BarColor = v.App.BarColor.Value()
	stackview.ItemTintColor = v.App.AllItemTintColor.Value()
	stackview.ItemTitleStyle, _ = v.App.AllItemTitleStyle.Value().(*text.Style)
	stackview.TitleStyle, _ = v.App.TitleStyle.Value().(*text.Style)
	return view.Model{
		Children: []view.View{stackview},
	}
}

type StackConfigureView struct {
	view.Embed
	App *StackApp
}

func NewStackConfigureView() *StackConfigureView {
	return &StackConfigureView{}
}

func (v *StackConfigureView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(&v.App.ItemTintColor)
		v.Subscribe(&v.App.ItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.App.ItemTintColor)
		v.Unsubscribe(&v.App.ItemTitleStyle)
	}
}

func (v *StackConfigureView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	button1 := view.NewButton()
	button1.String = "Toggle Bar Color"
	button1.OnPress = func() {
		if v.App.BarColor.Value() == nil {
			v.App.BarColor.SetValue(colornames.Red)
		} else {
			v.App.BarColor.SetValue(nil)
		}
	}
	l.Add(button1, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(50)
	})

	button4 := view.NewButton()
	button4.String = "Toggle Title Style"
	button4.OnPress = func() {
		if v.App.TitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Green)
			style.SetFont(text.DefaultBoldFont(20))
			v.App.TitleStyle.SetValue(style)
		} else {
			v.App.TitleStyle.SetValue(nil)
		}
	}
	l.Add(button4, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(50)
	})

	button2 := view.NewButton()
	button2.String = "Toggle All Item Color"
	button2.OnPress = func() {
		if v.App.AllItemTintColor.Value() == nil {
			v.App.AllItemTintColor.SetValue(colornames.Red)
		} else {
			v.App.AllItemTintColor.SetValue(nil)
		}
	}
	l.Add(button2, func(s *constraint.Solver) {
		s.Top(150)
		s.Left(50)
	})

	button3 := view.NewButton()
	button3.String = "Toggle All Item Style"
	button3.OnPress = func() {
		if v.App.AllItemTitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Green)
			style.SetFont(text.DefaultBoldFont(20))
			v.App.AllItemTitleStyle.SetValue(style)
		} else {
			v.App.AllItemTitleStyle.SetValue(nil)
		}
	}
	l.Add(button3, func(s *constraint.Solver) {
		s.Top(200)
		s.Left(50)
	})

	button5 := view.NewButton()
	button5.String = "Toggle Item Color"
	button5.OnPress = func() {
		if v.App.ItemTintColor.Value() == nil {
			v.App.ItemTintColor.SetValue(colornames.Orange)
		} else {
			v.App.ItemTintColor.SetValue(nil)
		}
	}
	l.Add(button5, func(s *constraint.Solver) {
		s.Top(250)
		s.Left(50)
	})

	button6 := view.NewButton()
	button6.String = "Toggle Item Style"
	button6.OnPress = func() {
		if v.App.ItemTitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Yellow)
			style.SetFont(text.DefaultBoldFont(20))
			v.App.ItemTitleStyle.SetValue(style)
		} else {
			v.App.ItemTitleStyle.SetValue(nil)
		}
	}
	l.Add(button6, func(s *constraint.Solver) {
		s.Top(300)
		s.Left(50)
	})

	leftItem := ios.NewTitleStackBarItem("TEST")
	leftItem.TintColor = v.App.ItemTintColor.Value()
	leftItem.TitleStyle, _ = v.App.ItemTitleStyle.Value().(*text.Style)
	leftItem.OnPress = func() {
		fmt.Println("Left Item on Press")
	}

	rightItem := ios.NewImageStackBarItem(application.MustLoadImage("checkbox_checked"))
	rightItem.OnPress = func() {
		fmt.Println("Right Item on Press")
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			&ios.StackBar{
				Title:      "Title",
				LeftItems:  []*ios.StackBarItem{leftItem},
				RightItems: []*ios.StackBarItem{rightItem},
			},
		},
	}
}

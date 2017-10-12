package ios

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewStackView", func() view.View {
		return NewStackView()
	})
}

type StackApp struct {
	IosStack          *ios.Stack
	AndroidStack      *android.Stack
	BarColor          comm.ColorValue
	ItemIconTint      comm.ColorValue
	ItemTitleStyle    comm.InterfaceValue
	TitleStyle        comm.InterfaceValue
	AllItemIconTint   comm.ColorValue
	AllItemTitleStyle comm.InterfaceValue
}

func NewStackView() view.View {
	app := &StackApp{
		IosStack:     &ios.Stack{},
		AndroidStack: &android.Stack{},
	}

	view1 := NewStackConfigureView()
	view1.App = app
	if runtime.GOOS == "android" {
		app.AndroidStack.SetViews(view1)
	} else {
		app.IosStack.SetViews(view1)
	}
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
		v.Subscribe(&v.App.AllItemIconTint)
		v.Subscribe(&v.App.AllItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.App.BarColor)
		v.Unsubscribe(&v.App.TitleStyle)
		v.Unsubscribe(&v.App.AllItemIconTint)
		v.Unsubscribe(&v.App.AllItemTitleStyle)
	}
}

func (v *StackAppView) Build(ctx view.Context) view.Model {
	var child view.View
	if runtime.GOOS == "android" {
		stackview := android.NewStackView()
		stackview.Stack = v.App.AndroidStack
		stackview.BarColor = v.App.BarColor.Value()
		stackview.ItemIconTint = v.App.AllItemIconTint.Value()
		stackview.ItemTitleStyle, _ = v.App.AllItemTitleStyle.Value().(*text.Style)
		stackview.TitleStyle, _ = v.App.TitleStyle.Value().(*text.Style)
		child = stackview
	} else {
		stackview := ios.NewStackView()
		stackview.Stack = v.App.IosStack
		stackview.BarColor = v.App.BarColor.Value()
		stackview.ItemIconTint = v.App.AllItemIconTint.Value()
		stackview.ItemTitleStyle, _ = v.App.AllItemTitleStyle.Value().(*text.Style)
		stackview.TitleStyle, _ = v.App.TitleStyle.Value().(*text.Style)
		child = stackview
	}
	return view.Model{
		Children: []view.View{child},
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
		v.Subscribe(&v.App.ItemIconTint)
		v.Subscribe(&v.App.ItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.App.ItemIconTint)
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
		if v.App.AllItemIconTint.Value() == nil {
			v.App.AllItemIconTint.SetValue(colornames.Red)
		} else {
			v.App.AllItemIconTint.SetValue(nil)
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
		if v.App.ItemIconTint.Value() == nil {
			v.App.ItemIconTint.SetValue(colornames.Orange)
		} else {
			v.App.ItemIconTint.SetValue(nil)
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

	leftIosItem := ios.NewStackBarItem()
	leftIosItem.Title = "TEST"
	leftIosItem.TitleStyle, _ = v.App.ItemTitleStyle.Value().(*text.Style)
	leftIosItem.OnPress = func() {
		fmt.Println("Left Item on Press")
	}

	rightIosItem := ios.NewStackBarItem()
	rightIosItem.Icon = application.MustLoadImage("checkbox_checked")
	rightIosItem.IconTint = v.App.ItemIconTint.Value()
	rightIosItem.OnPress = func() {
		fmt.Println("Right Item on Press")
	}

	leftAndroidItem := android.NewStackBarItem()
	if style, ok := v.App.ItemTitleStyle.Value().(*text.Style); ok {
		leftAndroidItem.StyledTitle = text.NewStyledText("TEST", style)
	} else {
		leftAndroidItem.Title = "TEST"
	}
	leftAndroidItem.OnPress = func() {
		fmt.Println("Left Item on Press")
	}

	rightAndroidItem := android.NewStackBarItem()
	rightAndroidItem.Icon = application.MustLoadImage("checkbox_checked")
	rightAndroidItem.IconTint = v.App.ItemIconTint.Value()
	rightAndroidItem.OnPress = func() {
		fmt.Println("Right Item on Press")
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			&ios.StackBar{
				Title:      "Title",
				LeftItems:  []*ios.StackBarItem{leftIosItem},
				RightItems: []*ios.StackBarItem{rightIosItem},
			},
			&android.StackBar{
				Title: "Title",
				Items: []*android.StackBarItem{leftAndroidItem, rightAndroidItem},
			},
		},
	}
}

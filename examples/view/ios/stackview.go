package ios

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/examples/internal"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewStackView", func() view.View {
		return NewStackAppView(NewStackApp())
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

func NewStackApp() *StackApp {
	app := &StackApp{
		IosStack:     &ios.Stack{},
		AndroidStack: &android.Stack{},
	}

	view1 := NewStackConfigureView(app)
	if runtime.GOOS == "android" {
		app.AndroidStack.SetViews(view1)
	} else {
		app.IosStack.SetViews(view1)
	}
	return app
}

type StackAppView struct {
	view.Embed
	app *StackApp
}

func NewStackAppView(app *StackApp) *StackAppView {
	return &StackAppView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *StackAppView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(&v.app.BarColor)
		v.Subscribe(&v.app.TitleStyle)
		v.Subscribe(&v.app.AllItemIconTint)
		v.Subscribe(&v.app.AllItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.app.BarColor)
		v.Unsubscribe(&v.app.TitleStyle)
		v.Unsubscribe(&v.app.AllItemIconTint)
		v.Unsubscribe(&v.app.AllItemTitleStyle)
	}
}

func (v *StackAppView) Build(ctx view.Context) view.Model {
	var child view.View
	if runtime.GOOS == "android" {
		stackview := android.NewStackView()
		stackview.Stack = v.app.AndroidStack
		stackview.BarColor = v.app.BarColor.Value()
		stackview.ItemIconTint = v.app.AllItemIconTint.Value()
		stackview.ItemTitleStyle, _ = v.app.AllItemTitleStyle.Value().(*text.Style)
		stackview.TitleStyle, _ = v.app.TitleStyle.Value().(*text.Style)
		child = stackview
	} else {
		stackview := ios.NewStackView()
		stackview.Stack = v.app.IosStack
		stackview.BarColor = v.app.BarColor.Value()
		stackview.ItemIconTint = v.app.AllItemIconTint.Value()
		stackview.ItemTitleStyle, _ = v.app.AllItemTitleStyle.Value().(*text.Style)
		stackview.TitleStyle, _ = v.app.TitleStyle.Value().(*text.Style)
		child = stackview
	}
	return view.Model{
		Children: []view.View{child},
	}
}

type StackConfigureView struct {
	view.Embed
	app *StackApp
}

func NewStackConfigureView(app *StackApp) *StackConfigureView {
	return &StackConfigureView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *StackConfigureView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(&v.app.ItemIconTint)
		v.Subscribe(&v.app.ItemTitleStyle)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(&v.app.ItemIconTint)
		v.Unsubscribe(&v.app.ItemTitleStyle)
	}
}

func (v *StackConfigureView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Toggle Bar Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	button := view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.BarColor.Value() == nil {
			v.app.BarColor.SetValue(colornames.Red)
		} else {
			v.app.BarColor.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle Title Style:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.TitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Green)
			style.SetFont(text.DefaultBoldFont(20))
			v.app.TitleStyle.SetValue(style)
		} else {
			v.app.TitleStyle.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle All Item Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.AllItemIconTint.Value() == nil {
			v.app.AllItemIconTint.SetValue(colornames.Red)
		} else {
			v.app.AllItemIconTint.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle All Item Style:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.AllItemTitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Green)
			style.SetFont(text.DefaultBoldFont(20))
			v.app.AllItemTitleStyle.SetValue(style)
		} else {
			v.app.AllItemTitleStyle.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle Item Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.ItemIconTint.Value() == nil {
			v.app.ItemIconTint.SetValue(colornames.Orange)
		} else {
			v.app.ItemIconTint.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle Item Style:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.ItemTitleStyle.Value() == nil {
			style := &text.Style{}
			style.SetTextColor(colornames.Yellow)
			style.SetFont(text.DefaultBoldFont(20))
			v.app.ItemTitleStyle.SetValue(style)
		} else {
			v.app.ItemTitleStyle.SetValue(nil)
		}
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textIosItem := ios.NewStackBarItem()
	textIosItem.Title = "Back"
	textIosItem.TitleStyle, _ = v.app.ItemTitleStyle.Value().(*text.Style)
	textIosItem.OnPress = func() {
		fmt.Println("Left Item on Press")
		if internal.BackRelay != nil {
			internal.BackRelay.Signal()
		}
	}

	imageIosItem := ios.NewStackBarItem()
	imageIosItem.Icon = application.MustLoadImage("checkbox_checked")
	imageIosItem.IconTint = v.app.ItemIconTint.Value()
	imageIosItem.OnPress = func() {
		fmt.Println("Right Item on Press")
	}

	leftAndroidItem := android.NewStackBarItem()
	if style, ok := v.app.ItemTitleStyle.Value().(*text.Style); ok {
		leftAndroidItem.StyledTitle = text.NewStyledText("Back", style)
	} else {
		leftAndroidItem.Title = "Back"
	}
	leftAndroidItem.OnPress = func() {
		fmt.Println("Left Item on Press")
		if internal.BackRelay != nil {
			internal.BackRelay.Signal()
		}
	}

	rightAndroidItem := android.NewStackBarItem()
	rightAndroidItem.Icon = application.MustLoadImage("checkbox_checked")
	rightAndroidItem.IconTint = v.app.ItemIconTint.Value()
	rightAndroidItem.OnPress = func() {
		fmt.Println("Right Item on Press")
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Options: []view.Option{
			&ios.StackBar{
				Title:      "Title",
				LeftItems:  []*ios.StackBarItem{textIosItem},
				RightItems: []*ios.StackBarItem{imageIosItem},
			},
			&android.StackBar{
				Title: "Title",
				Items: []*android.StackBarItem{leftAndroidItem, rightAndroidItem},
			},
		},
	}
}

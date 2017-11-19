package examples

import (
	"runtime"

	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	applicationex "gomatcha.io/matcha/examples/application"
	bridgeex "gomatcha.io/matcha/examples/bridge"
	"gomatcha.io/matcha/examples/customview"
	"gomatcha.io/matcha/examples/insta"
	"gomatcha.io/matcha/examples/internal"
	layoutex "gomatcha.io/matcha/examples/layout"
	paintex "gomatcha.io/matcha/examples/paint"
	"gomatcha.io/matcha/examples/settings"
	"gomatcha.io/matcha/examples/todo"
	viewex "gomatcha.io/matcha/examples/view"
	androidex "gomatcha.io/matcha/examples/view/android"
	iosex "gomatcha.io/matcha/examples/view/ios"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples NewExamplesView", func() view.View {
		app := NewApp()
		return NewAppView(app)
	})
}

// Wrapper for the ios.Stack and android.Stack
type Stack interface {
	SetViews(vs ...view.View)
	Views() []view.View
	Push(vs view.View)
	Pop()
}

type Example struct {
	Title       string
	PushOnStack bool
	Description string
	View        view.View
}

type App struct {
	Stack Stack

	ChildRelay *comm.Relay
	Child      view.View
}

func NewApp() *App {
	app := &App{
		ChildRelay: &comm.Relay{},
	}
	if runtime.GOOS == "android" {
		app.Stack = &android.Stack{}
	} else {
		app.Stack = &ios.Stack{}
	}
	v := NewExamplesView(app)
	app.Stack.SetViews(v)
	return app
}

type AppView struct {
	view.Embed
	app     *App
	backKey comm.Id
}

func NewAppView(app *App) *AppView {
	return &AppView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *AppView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		internal.ExamplesRelay = &comm.Relay{}
		v.backKey = internal.ExamplesRelay.Notify(func() {
			// Pop to root and clear the child
			v.app.Stack.Pop()
			v.app.Child = nil
			v.app.ChildRelay.Signal()
		})
		v.Subscribe(v.app.ChildRelay)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		internal.ExamplesRelay.Unnotify(v.backKey)
		internal.ExamplesRelay = nil
		v.Unsubscribe(v.app.ChildRelay)
	}
}

func (v *AppView) Build(ctx view.Context) view.Model {
	// If user has selected an example, display it.
	if v.app.Child != nil {
		return view.Model{Children: []view.View{v.app.Child}}
	}

	// Otherwise display the stack view
	var stack view.View
	if runtime.GOOS == "android" {
		stackview := android.NewStackView()
		stackview.Stack = v.app.Stack.(*android.Stack)
		stack = stackview
	} else {
		stackview := ios.NewStackView()
		stackview.Stack = v.app.Stack.(*ios.Stack)
		stack = stackview
	}
	return view.Model{Children: []view.View{stack}}
}

type ExamplesView struct {
	view.Embed
	app *App
}

func NewExamplesView(app *App) *ExamplesView {
	return &ExamplesView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *ExamplesView) Build(ctx view.Context) view.Model {
	childLayouter := &table.Layouter{StartEdge: layout.EdgeTop}

	// Sections
	sections := map[string][]*Example{
		"Examples": {
			{"Settings", false, "Example of a settings app.\n\n\ngomatcha.io/matcha/examples/settings", settings.NewRootView()},
			{"Instagram", false, "Example of an photo-sharing app.\n\n\ngomatcha.io/matcha/examples/insta", insta.NewRootView()},
			{"Todo App", false, "Example of a basic todo app.\n\n\ngomatcha.io/matcha/examples/todo", todo.NewRootView()},
		},
		"General": {
			// {"animate.NewView", "", animate.NewView()},
			{"Painters", true, "Example of various paint options.\n\n\ngomatcha.io/matcha/examples/paint NewPaintView", paintex.NewPaintView()},
			{"Constraint Layout", true, "Example of a complex layout using constraints. \n\n\ngomatcha.io/matcha/examples/layout NewConstraintsView", layoutex.NewConstraintsView()},
			{"Table Layout", true, "Example of a vertical (bottom to top) table and a horizontal (left to right) table layout. \n\n\ngomatcha.io/matcha/examples/layout NewTableView", layoutex.NewTableView()},
			{"Native Bridge", true, "Example of how to call native functions from Go, and Go functions from native code. See gomatcha.io/guide/native-bridge/ for more details. \n\n\ngomatcha.io/matcha/examples/bridge NewBridgeView", bridgeex.NewBridgeView()},
			{"Custom Views", true, "Example of creating a custom view, that displays the user's camera. \n\n\ngomatcha.io/matcha/examples/customview NewView", customview.NewView()},
		},
		"Views": {
			{"Alerts", true, "Example of how to display alerts. \n\n\ngomatcha.io/matcha/examples/view NewAlertView", viewex.NewAlertView()},
			{"Button", true, "Example of an enabled button and a disabled button. \n\n\ngomatcha.io/matcha/examples/view NewButtonView", viewex.NewButtonView()},
			{"Image View", true, "Example of various image view properties. \n\n\ngomatcha.io/matcha/examples/view NewImageView", viewex.NewImageView()},
			{"Scroll View", true, "Example of a scroll view with a label that tracks the scroll position and a button that scrolls to a point. \n\n\ngomatcha.io/matcha/examples/view NewScrollView", viewex.NewScrollView()},
			{"Slider", true, "Example of an enabled sider and a disabled slider that tracks the enabled slider. \n\n\ngomatcha.io/matcha/examples/view NewSliderView", viewex.NewSliderView()},
			{"Switch", true, "Example of an enabled switch and a disabled switch that tracks the enabled switch. \n\n\ngomatcha.io/matcha/examples/view NewSwitchView", viewex.NewSwitchView()},
			{"Text View", true, "Example of various text view properties. \n\n\ngomatcha.io/matcha/examples/view NewTextView", viewex.NewTextView()},
			// Text input example
		},
		"iOS": {
			{"Activity Indicator", true, "Example of how to show/hide the activity indicator. \n\n\ngomatcha.io/matcha/examples/view/ios NewActivityIndicatorView", iosex.NewActivityIndicatorView()},
			{"Segment View", true, "Example of various segment view properties. \n\n\ngomatcha.io/matcha/examples/view/ios NewSegmentView", iosex.NewSegmentView()},
			{"Status Bar", true, "Example of how to toggle the status bar style and color. \n\n\ngomatcha.io/matcha/examples/view/ios NewStatusBarView", iosex.NewStatusBarView()},
			{"Stack View", false, "Example of various stack view properties. \n\n\ngomatcha.io/matcha/examples/view/ios NewStackView", iosex.NewStackAppView(iosex.NewStackApp())},
			{"Tab View", true, "Example of how to use and customize a tab view. \n\n\ngomatcha.io/matcha/examples/view/ios NewTabView", iosex.NewTabView(iosex.NewTabApp())},
			{"Progress View", true, "Example of how to use a progress view. \n\n\ngomatcha.io/matcha/examples/view/ios NewProgressView", iosex.NewProgressView()},
			// {"Navigation", "\n\n\ngomatcha.io/matcha/examples/view/ios NewNavigationView", iosex.NewNavigationView()},
		},
		"Android": {
			{"Pager View", true, "\n\n\ngomatcha.io/matcha/examples/view/android NewPagerView", androidex.NewPagerView()},
			{"Stack View", false, "\n\n\ngomatcha.io/matcha/examples/view/android NewStackView", androidex.NewStackView()},
			{"Status Bar", true, "\n\n\ngomatcha.io/matcha/examples/view/android NewStatusBarView", androidex.NewStatusBarView()},
		},
		"Miscellaneous": {
			{"Device Orientation", true, "Example of how to get the current device orientation and listen to orientation changes.\n\n\ngomatcha.io/matcha/examples/application NewOrientationView", applicationex.NewOrientationView()},
			// Shake device
			{"Adding/Removing Views", true, "\n\n\ngomatcha.io/matcha/examples/view NewAddRemoveView", viewex.NewAddRemoveView()},
		},
	}

	for _, i := range []string{"Examples", "General", "Views", "iOS", "Android", "Miscellaneous"} {
		// Create header for section.
		header := settings.NewSpacerHeader()
		header.Title = i
		childLayouter.Add(header, nil)

		// Create example items for section.
		items := []view.View{}
		for _, j := range sections[i] {
			example := j
			item := settings.NewBasicCell()
			item.Title = j.Title
			item.OnTap = func() {
				v.app.Stack.Push(NewExamplesDetailView(v.app, example))
			}
			item.Chevron = true
			items = append(items, item)
		}

		// Add separators around items.
		for _, j := range settings.AddSeparators(items, 30) {
			childLayouter.Add(j, nil)
		}
	}

	// Add footer.
	footer := settings.NewSpacer()
	footer.Height = 50
	childLayouter.Add(footer, nil)

	sv := view.NewScrollView()
	sv.ContentPainter = &paint.Style{BackgroundColor: settings.BackgroundColor}
	sv.ContentLayouter = childLayouter
	sv.ContentChildren = childLayouter.Views()

	return view.Model{
		Children: []view.View{sv},
		Painter:  &paint.Style{BackgroundColor: settings.BackgroundColor},
		Options: []view.Option{
			&ios.StackBar{Title: "Matcha"},
			&android.StackBar{Title: "Matcha"},
		},
	}
}

type ExamplesDetailView struct {
	view.Embed
	app     *App
	example *Example
}

func NewExamplesDetailView(app *App, e *Example) *ExamplesDetailView {
	return &ExamplesDetailView{
		Embed:   view.NewEmbed(app),
		app:     app,
		example: e,
	}
}

func (v *ExamplesDetailView) Build(ctx view.Context) view.Model {
	childLayouter := &table.Layouter{StartEdge: layout.EdgeTop}

	childLayouter.Add(settings.NewSpacer(), nil)

	item := settings.NewLargeCell()
	item.Title = v.example.Title
	item.OnTap = func() {
		if v.example.PushOnStack {
			v.app.Stack.Push(v.example.View)
		} else {
			v.app.Child = v.example.View
			v.app.ChildRelay.Signal()
		}
	}
	for _, j := range settings.AddSeparators([]view.View{item}, 30) {
		childLayouter.Add(j, nil)
	}

	header := settings.NewSpacerDescription()
	header.Description = v.example.Description
	childLayouter.Add(header, nil)

	sv := view.NewScrollView()
	sv.ContentPainter = &paint.Style{BackgroundColor: settings.BackgroundColor}
	sv.ContentLayouter = childLayouter
	sv.ContentChildren = childLayouter.Views()

	return view.Model{
		Children: []view.View{sv},
		Painter:  &paint.Style{BackgroundColor: settings.BackgroundColor},
		Options: []view.Option{
			&ios.StackBar{Title: "Matcha"},
			&android.StackBar{Title: "Matcha"},
		},
	}
}

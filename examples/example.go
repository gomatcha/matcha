package examples

import (
	"runtime"

	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	applicationex "gomatcha.io/matcha/examples/application"
	bridgeex "gomatcha.io/matcha/examples/bridge"
	"gomatcha.io/matcha/examples/customview"
	"gomatcha.io/matcha/examples/insta"
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
		return NewRootView(app)
	})
}

type Section struct {
	Title    string
	Examples []*Example
}

type Example struct {
	Title       string
	Description string
	View        view.View
}

type App struct {
	IosStack     *ios.Stack
	AndroidStack *android.Stack
	Child        view.View
	ChildRelay   *comm.Relay
}

func NewApp() *App {
	app := &App{
		AndroidStack: &android.Stack{},
		IosStack:     &ios.Stack{},
		ChildRelay:   &comm.Relay{},
	}
	initialView := NewExamplesView(app)
	app.AndroidStack.SetViews(initialView)
	app.IosStack.SetViews(initialView)
	return app
}

type RootView struct {
	view.Embed
	app      *App
	shakeKey comm.Id
}

func NewRootView(app *App) *RootView {
	return &RootView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *RootView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.shakeKey = application.ShakeNotifier.Notify(func() {
			// Pop to root and clear the child, if user shakes the device.
			v.app.AndroidStack.Pop()
			v.app.IosStack.Pop()
			v.app.Child = nil
			v.app.ChildRelay.Signal()
		})
		v.Subscribe(v.app.ChildRelay)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		application.ShakeNotifier.Unnotify(v.shakeKey)
		v.Unsubscribe(v.app.ChildRelay)
	}
}

func (v *RootView) Build(ctx view.Context) view.Model {
	// If user has selected an example, display it.
	if v.app.Child != nil {
		return view.Model{Children: []view.View{v.app.Child}}
	}

	// Otherwise display the stack view
	var stack view.View
	if runtime.GOOS == "android" {
		stackview := android.NewStackView()
		stackview.Stack = v.app.AndroidStack
		stack = stackview
	} else {
		stackview := ios.NewStackView()
		stackview.Stack = v.app.IosStack
		stack = stackview
	}
	return view.Model{Children: []view.View{stack}}
}

type ExamplesView struct {
	view.Embed
	app      *App
	sections []Section
}

func NewExamplesView(app *App) *ExamplesView {
	sections := []Section{
		{
			Title: "Examples",
			Examples: []*Example{
				&Example{"Settings", "Example of a settings app.\n\n\ngomatcha.io/matcha/examples/settings", settings.NewRootView()},
				&Example{"Instagram", "Example of an photo-sharing app.\n\n\ngomatcha.io/matcha/examples/insta", insta.NewRootView()},
				&Example{"Todo App", "Example of a basic todo app.\n\n\ngomatcha.io/matcha/examples/todo", todo.NewRootView()},
			},
		},
		{
			Title: "General",
			Examples: []*Example{
				// {"animate.NewView", "", animate.NewView()},
				&Example{"Painters", "Example of various paint options.\n\n\ngomatcha.io/matcha/examples/paint NewPaintView", paintex.NewPaintView()},
				&Example{"Constraint Layout", "Example of a complex layout using constraints. \n\n\ngomatcha.io/matcha/examples/layout NewConstraintsView", layoutex.NewConstraintsView()},
				&Example{"Table Layout", "Example of a vertical (bottom to top) table and a horizontal (left to right) table layout. \n\n\ngomatcha.io/matcha/examples/layout NewTableView", layoutex.NewTableView()},
				&Example{"Native Bridge", "Example of how to call native functions from Go, and Go functions from native code. See gomatcha.io/guide/native-bridge/ for more details. \n\n\ngomatcha.io/matcha/examples/bridge NewBridgeView", bridgeex.NewBridgeView()},
				&Example{"Custom Views", "Example of creating a custom view, that displays the user's camera. \n\n\ngomatcha.io/matcha/examples/customview NewView", customview.NewView()},
			},
		},
		{
			Title: "Views",
			Examples: []*Example{
				&Example{"Alerts", "Example of how to display alerts. \n\n\ngomatcha.io/matcha/examples/view NewAlertView", viewex.NewAlertView()},
				&Example{"Button", "Example of an enabled button and a disabled button. \n\n\ngomatcha.io/matcha/examples/view NewButtonView", viewex.NewButtonView()},
				&Example{"Image View", "Example of various image view properties. \n\n\ngomatcha.io/matcha/examples/view NewImageView", viewex.NewImageView()},
				&Example{"Scroll View", "Example of a scroll view with a label that tracks the scroll position and a button that scrolls to a point. \n\n\ngomatcha.io/matcha/examples/view NewScrollView", viewex.NewScrollView()},
				&Example{"Slider", "Example of an enabled sider and a disabled slider that tracks the enabled slider. \n\n\ngomatcha.io/matcha/examples/view NewSliderView", viewex.NewSliderView()},
				&Example{"Switch", "Example of an enabled switch and a disabled switch that tracks the enabled switch. \n\n\ngomatcha.io/matcha/examples/view NewSwitchView", viewex.NewSwitchView()},
				&Example{"Text View", "Example of various text view properties. \n\n\ngomatcha.io/matcha/examples/view NewTextView", viewex.NewTextView()},
				// Text input example
			},
		},
		{
			Title: "iOS",
			Examples: []*Example{
				&Example{"Activity Indicator", "Example of how to show/hide the activity indicator. \n\n\ngomatcha.io/matcha/examples/view/ios NewActivityIndicatorView", iosex.NewActivityIndicatorView()},
				&Example{"Segment View", "Example of various segment view properties. \n\n\ngomatcha.io/matcha/examples/view/ios NewSegmentView", iosex.NewSegmentView()},
				&Example{"Status Bar", "Example of how to toggle the status bar style and color. \n\n\ngomatcha.io/matcha/examples/view/ios NewStatusBarView", iosex.NewStatusBarView()},
				&Example{"Stack View", "Example of various stack view properties. \n\n\ngomatcha.io/matcha/examples/view/ios NewStackView", iosex.NewStackAppView(iosex.NewStackApp())},
				&Example{"Tab View", "Example of how to use and customize a tab view. \n\n\ngomatcha.io/matcha/examples/view/ios NewTabView", iosex.NewTabView(iosex.NewTabApp())},
				&Example{"Progress View", "Example of how to use a progress view. \n\n\ngomatcha.io/matcha/examples/view/ios NewProgressView", iosex.NewProgressView()},
				// &Example{"Navigation", "\n\n\ngomatcha.io/matcha/examples/view/ios NewNavigationView", iosex.NewNavigationView()},
			},
		},
		{
			Title: "Android",
			Examples: []*Example{
				&Example{"Pager View", "\n\n\ngomatcha.io/matcha/examples/view/android NewPagerView", androidex.NewPagerView()},
				&Example{"Stack View", "\n\n\ngomatcha.io/matcha/examples/view/android NewStackView", androidex.NewStackView()},
				&Example{"Status Bar", "\n\n\ngomatcha.io/matcha/examples/view/android NewStatusBarView", androidex.NewStatusBarView()},
			},
		},
		{
			Title: "Miscellaneous",
			Examples: []*Example{
				&Example{"Device Orientation", "Example of how to get the current device orientation and listen to orientation changes.\n\n\ngomatcha.io/matcha/examples/application NewOrientationView", applicationex.NewOrientationView()},
				&Example{"Adding/Removing Views", "\n\n\ngomatcha.io/matcha/examples/view NewAddRemoveView", viewex.NewAddRemoveView()},
			},
		},
	}

	return &ExamplesView{
		Embed:    view.NewEmbed(app),
		app:      app,
		sections: sections,
	}
}

func (v *ExamplesView) Build(ctx view.Context) view.Model {
	childLayouter := &table.Layouter{StartEdge: layout.EdgeTop}

	// Add header.
	header := settings.NewSpacerDescription()
	header.Description = "Shake the device to return back to this list."
	childLayouter.Add(header, nil)

	for _, i := range v.sections {
		// Create header for section.
		header := settings.NewSpacerHeader()
		header.Title = i.Title
		childLayouter.Add(header, nil)

		// Create example items for section.
		items := []view.View{}
		for _, j := range i.Examples {
			example := j
			item := settings.NewBasicCell()
			item.Title = j.Title
			item.OnTap = func() {
				detailView := NewExamplesDetailView(v.app, example)
				v.app.IosStack.Push(detailView)
				v.app.AndroidStack.Push(detailView)
			}
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
		v.app.Child = v.example.View
		v.app.ChildRelay.Signal()
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

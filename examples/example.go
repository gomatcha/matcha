package examples

import (
	"github.com/gomatcha/matcha/application"
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/comm"
	applicationex "github.com/gomatcha/matcha/examples/application"
	bridgeex "github.com/gomatcha/matcha/examples/bridge"
	"github.com/gomatcha/matcha/examples/customview"
	"github.com/gomatcha/matcha/examples/insta"
	layoutex "github.com/gomatcha/matcha/examples/layout"
	paintex "github.com/gomatcha/matcha/examples/paint"
	"github.com/gomatcha/matcha/examples/settings"
	"github.com/gomatcha/matcha/examples/todo"
	viewex "github.com/gomatcha/matcha/examples/view"
	"github.com/gomatcha/matcha/examples/view/android"
	"github.com/gomatcha/matcha/examples/view/ios"
	"github.com/gomatcha/matcha/layout"
	"github.com/gomatcha/matcha/layout/table"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples NewExamplesView", func() view.View {
		return NewExamplesView()
	})
}

type Section struct {
	Title    string
	Examples []Example
}

type Example struct {
	Title       string
	Description string
	View        view.View
}

type ExamplesView struct {
	view.Embed
	Sections []Section
	child    view.View
	shakeKey comm.Id
}

func NewExamplesView() *ExamplesView {
	sections := []Section{
		{
			Title: "Examples",
			Examples: []Example{
				{"Settings", "", settings.NewRootView()},
				{"Instagram", "", insta.NewRootView()},
				{"Todo App", "", todo.NewRootView()},
			},
		},
		{
			Title: "General",
			Examples: []Example{
				// {"animate.NewView", "", animate.NewView()},
				{"Device Orientation", "", applicationex.NewOrientationView()},
				{"Native Bridge", "", bridgeex.NewBridgeView()},
				// {"complex.NewNestedView", "", complex.NewNestedView()},
				{"Custom Views", "", customview.NewView()},
				{"Constraints Layout", "", layoutex.NewConstraintsView()},
				{"Table Layout", "", layoutex.NewTableView()},
				{"Painters", "", paintex.NewPaintView()},
				{"Adding/Removing Views", "", viewex.NewAddRemoveView()},
			},
		},
		{
			Title: "Views",
			Examples: []Example{
				{"Alerts", "", viewex.NewAlertView()},
				{"Button", "", viewex.NewButtonView()},
				{"Image View", "", viewex.NewImageView()},
				{"Scroll View", "", viewex.NewScrollView()},
				{"Slider", "", viewex.NewSliderView()},
				{"Switch View", "", viewex.NewSwitchView()},
				{"Text View", "", viewex.NewTextView()},
			},
		},
		{
			Title: "iOS",
			Examples: []Example{
				{"Activity Indicator", "", ios.NewActivityIndicatorView()},
				{"Navigation", "", ios.NewNavigationView()},
				{"Segment View", "", ios.NewSegmentView()},
				{"Stack View", "", ios.NewStackView()},
				{"Status Bar", "", ios.NewStatusBarView()},
				{"Tab View", "", ios.NewTabView()},
				{"Progress View", "", ios.NewProgressView()},
			},
		},
		{
			Title: "Android",
			Examples: []Example{
				{"Pager View", "", android.NewPagerView()},
				{"Stack View", "", android.NewStackView()},
				{"Status Bar", "", android.NewStatusBarView()},
			},
		},
	}

	return &ExamplesView{
		Sections: sections,
	}
}

func (v *ExamplesView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.shakeKey = application.ShakeNotifier.Notify(func() {
			v.child = nil
			v.Signal()
		})
	} else if view.ExitsStage(from, to, view.StageMounted) {
		application.ShakeNotifier.Unnotify(v.shakeKey)
	}
}

func (v *ExamplesView) Build(ctx view.Context) view.Model {
	// If user has selected an example, display it.
	if v.child != nil {
		return view.Model{Children: []view.View{v.child}}
	}

	childLayouter := &table.Layouter{StartEdge: layout.EdgeTop}

	// Add header.
	header := settings.NewSpacerDescription()
	header.Description = "Shake the device to return back to this list."
	childLayouter.Add(header, nil)

	for _, i := range v.Sections {
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
				v.child = example.View
				v.Signal()
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
	}
}

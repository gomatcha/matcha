package examples

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/examples/application"
	bridgeex "gomatcha.io/matcha/examples/bridge"
	"gomatcha.io/matcha/examples/customview"
	"gomatcha.io/matcha/examples/insta"
	layoutex "gomatcha.io/matcha/examples/layout"
	paintex "gomatcha.io/matcha/examples/paint"
	"gomatcha.io/matcha/examples/settings"
	"gomatcha.io/matcha/examples/todo"
	viewex "gomatcha.io/matcha/examples/view"
	"gomatcha.io/matcha/examples/view/android"
	"gomatcha.io/matcha/examples/view/ios"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewExamplesView", func() view.View {
		return NewExamplesView()
	})
}

type ExamplesView struct {
	view.Embed
	child view.View
}

func NewExamplesView() *ExamplesView {
	return &ExamplesView{}
}

func (v *ExamplesView) Build(ctx view.Context) view.Model {
	if v.child != nil {
		return view.Model{Children: []view.View{v.child}}
	}

	childLayouter := &table.Layouter{
		StartEdge: layout.EdgeTop,
	}

	items := []struct {
		title string
		view  view.View
	}{
		// {"animate.NewView", animate.NewView()},
		{"application.NewOrientationView", application.NewOrientationView()},
		{"bridge.NewBridgeView", bridgeex.NewBridgeView()},
		// {"complex.NewNestedView", complex.NewNestedView()},
		{"customview.NewView", customview.NewView()},
		{"insta.NewRootView", insta.NewRootView()},
		{"layout.NewConstraintsView", layoutex.NewConstraintsView()},
		{"layout.NewTableView", layoutex.NewTableView()},
		{"paint.NewPaintView", paintex.NewPaintView()},
		{"settings.NewRootView", settings.NewRootView()},
		{"todo.NewRootView", todo.NewRootView()},
		{"view/android.NewPagerView", android.NewPagerView()},
		{"view/android.NewStackView", android.NewStackView()},
		{"view/android.NewStatusBarView", android.NewStatusBarView()},
		{"view/ios.NewActivityIndicatorView", ios.NewActivityIndicatorView()},
		{"view/ios.NewNavigationView", ios.NewNavigationView()},
		{"view/ios.NewSegmentView", ios.NewSegmentView()},
		{"view/ios.NewStackView", ios.NewStackView()},
		{"view/ios.NewStatusBarView", ios.NewStatusBarView()},
		{"view/ios.NewTabView", ios.NewTabView()},
		{"view.NewAddRemoveView", viewex.NewAddRemoveView()},
		{"view.NewAlertView", viewex.NewAlertView()},
		{"view.NewButtonView", viewex.NewButtonView()},
		{"view.NewImageView", viewex.NewImageView()},
		{"view.NewProgressView", viewex.NewProgressView()},
		{"view.NewScrollView", viewex.NewScrollView()},
		{"view.NewSliderView", viewex.NewSliderView()},
		{"view.NewSwitchView", viewex.NewSwitchView()},
		{"view.NewTextView", viewex.NewTextView()},
		{"view.NewUnknownView", viewex.NewUnknownView()},
	}

	for _, i := range items {
		item := i
		childView := NewExampleCell()
		childView.String = i.title
		childView.OnPress = func() {
			v.child = item.view
			v.Signal()
		}
		childLayouter.Add(childView, nil)
	}

	sv := view.NewScrollView()
	sv.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	sv.ContentLayouter = childLayouter
	sv.ContentChildren = childLayouter.Views()
	sv.PaintStyle = &paint.Style{BackgroundColor: colornames.Cyan}

	return view.Model{
		Children: []view.View{sv},
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

type ExampleCell struct {
	view.Embed
	String  string
	OnPress func()
}

func NewExampleCell() *ExampleCell {
	return &ExampleCell{}
}

func (v *ExampleCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
	})

	textView := view.NewTextView()
	textView.String = v.String
	textView.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	l.Add(textView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.CenterYEqual(l.CenterY())
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
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			pointer.GestureList{button},
		},
	}
}

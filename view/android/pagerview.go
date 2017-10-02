package android

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout/constraint"
	pbandroid "gomatcha.io/matcha/proto/view/android"
	"gomatcha.io/matcha/view"
)

// Pages represents a list of views to be shown in the PagerView. It can be manipulated outside of a Build() call.
type Pages struct {
	relay         comm.Relay
	children      []view.View
	selectedIndex int
}

// SetViews sets the child views displayed in the tabview.
func (s *Pages) SetViews(ss ...view.View) {
	s.children = ss
	s.relay.Signal()
}

// Views returns the child views displayed in the tabview.
func (s *Pages) Views() []view.View {
	return s.children
}

// SetSelectedIndex selects the tab at idx.
func (s *Pages) SetSelectedIndex(idx int) {
	if idx != s.selectedIndex {
		s.selectedIndex = idx
		s.relay.Signal()
	}
}

// SelectedIndex returns the index of the selected tab.
func (s *Pages) SelectedIndex() int {
	return s.selectedIndex
}

func (s *Pages) SelectedView() view.View {
	if s.selectedIndex > len(s.children)-1 {
		return nil
	}
	return s.children[s.selectedIndex]
}

// Notify implements the comm.Notifier interface.
func (s *Pages) Notify(f func()) comm.Id {
	return s.relay.Notify(f)
}

// Unnotify implements the comm.Notifier interface.
func (s *Pages) Unnotify(id comm.Id) {
	s.relay.Unnotify(id)
}

type PagerView struct {
	view.Embed
	Pages *Pages
}

// NewPagerView returns a new view.
// ViewPager and PagerTabStrip.
func NewPagerView() *PagerView {
	return &PagerView{
		Pages: &Pages{},
	}
}

// Lifecyle implements the view.View interface.
func (v *PagerView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		if v.Pages == nil {
			v.Pages = &Pages{}
		}
		v.Subscribe(v.Pages)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.Pages)
	}
}

func (v *PagerView) Update(v2 view.View) {
	v.Unsubscribe(v.Pages)

	view.CopyFields(v, v2)

	if v.Pages == nil {
		v.Pages = &Pages{}
	}
	v.Subscribe(v.Pages)
}

// Build implements the view.View interface.
func (v *PagerView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	childrenPb := []*pbandroid.PagerChildView{}
	for _, chld := range v.Pages.Views() {
		// Find the button
		var button *PagerButton

		for _, opts := range chld.Build(nil).Options {
			var ok bool
			if button, ok = opts.(*PagerButton); ok {
				break
			}
		}

		if button == nil {
			button = &PagerButton{
				Title: "Title",
			}
		}

		// Add the child.
		l.Add(chld, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.WidthEqual(l.MaxGuide().Width())
			s.HeightEqual(l.MaxGuide().Height())
		})

		// Add to protobuf.
		childrenPb = append(childrenPb, &pbandroid.PagerChildView{
			Title: button.Title,
			// Icon:         app.ImageMarshalProtobuf(button.Icon),
			// SelectedIcon: app.ImageMarshalProtobuf(button.SelectedIcon),
			// Badge:        button.Badge,
		})
	}

	// var selectedTextStyle *pbtext.TextStyle
	// if v.SelectedTextStyle != nil {
	// 	selectedTextStyle = v.SelectedTextStyle.MarshalProtobuf()
	// }

	// var unselectedTextStyle *pbtext.TextStyle
	// if v.UnselectedTextStyle != nil {
	// 	unselectedTextStyle = v.UnselectedTextStyle.MarshalProtobuf()
	// }

	return view.Model{
		Children:       l.Views(),
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/android PagerView",
		NativeViewState: internal.MarshalProtobuf(&pbandroid.PagerView{
			ChildViews:    childrenPb,
			SelectedIndex: int64(v.Pages.SelectedIndex()),
			// BarColor:            pb.ColorEncode(v.BarColor),
			// SelectedColor:       pb.ColorEncode(v.SelectedColor),
			// UnselectedColor:     pb.ColorEncode(v.UnselectedColor),
			// SelectedTextStyle:   selectedTextStyle,
			// UnselectedTextStyle: unselectedTextStyle,
		}),
		NativeFuncs: map[string]interface{}{
			"OnSelect": func(index int64) {
				v.Pages.SetSelectedIndex(int(index))
			},
		},
	}
}

type PagerButton struct {
	Title string
}

func (t *PagerButton) OptionKey() string {
	return "gomatcha.io/view/android PagerButton"
}

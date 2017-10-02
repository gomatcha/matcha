package ios

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout/constraint"
	pb "gomatcha.io/matcha/proto"
	pbtext "gomatcha.io/matcha/proto/text"
	pbios "gomatcha.io/matcha/proto/view/ios"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

// Tabs represents a list of views to be shown in the TabView. It can be manipulated outside of a Build() call.
type Tabs struct {
	relay         comm.Relay
	children      []view.View
	selectedIndex int
}

// SetViews sets the child views displayed in the tabview.
func (s *Tabs) SetViews(ss ...view.View) {
	s.children = ss
	s.relay.Signal()
}

// Views returns the child views displayed in the tabview.
func (s *Tabs) Views() []view.View {
	return s.children
}

// SetSelectedIndex selects the tab at idx.
func (s *Tabs) SetSelectedIndex(idx int) {
	if idx != s.selectedIndex {
		s.selectedIndex = idx
		s.relay.Signal()
	}
}

// SelectedIndex returns the index of the selected tab.
func (s *Tabs) SelectedIndex() int {
	return s.selectedIndex
}

func (s *Tabs) SelectedView() view.View {
	if s.selectedIndex > len(s.children)-1 {
		return nil
	}
	return s.children[s.selectedIndex]
}

// Notify implements the comm.Notifier interface.
func (s *Tabs) Notify(f func()) comm.Id {
	return s.relay.Notify(f)
}

// Unnotify implements the comm.Notifier interface.
func (s *Tabs) Unnotify(id comm.Id) {
	s.relay.Unnotify(id)
}

type TabView struct {
	view.Embed
	Tabs                *Tabs
	BarColor            color.Color
	SelectedTextStyle   *text.Style
	UnselectedTextStyle *text.Style
	SelectedColor       color.Color
	UnselectedColor     color.Color
}

// NewTabView returns a new view.
func NewTabView() *TabView {
	return &TabView{
		Tabs: &Tabs{},
	}
}

// Lifecyle implements the view.View interface.
func (v *TabView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		if v.Tabs == nil {
			v.Tabs = &Tabs{}
		}
		v.Subscribe(v.Tabs)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.Tabs)
	}
}

// Update implements the view.View interface.
func (v *TabView) Update(v2 view.View) {
	v.Unsubscribe(v.Tabs)

	view.CopyFields(v, v2)

	if v.Tabs == nil {
		v.Tabs = &Tabs{}
	}
	v.Subscribe(v.Tabs)
}

// Build implements the view.View interface.
func (v *TabView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	childrenPb := []*pbios.TabChildView{}
	for _, chld := range v.Tabs.Views() {
		// Find the button
		var button *TabButton

		for _, opts := range chld.Build(nil).Options {
			var ok bool
			if button, ok = opts.(*TabButton); ok {
				break
			}
		}

		if button == nil {
			button = &TabButton{
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
		childrenPb = append(childrenPb, &pbios.TabChildView{
			Title:        button.Title,
			Icon:         internal.ImageMarshalProtobuf(button.Icon),
			SelectedIcon: internal.ImageMarshalProtobuf(button.SelectedIcon),
			Badge:        button.Badge,
		})
	}

	var selectedTextStyle *pbtext.TextStyle
	if v.SelectedTextStyle != nil {
		selectedTextStyle = v.SelectedTextStyle.MarshalProtobuf()
	}

	var unselectedTextStyle *pbtext.TextStyle
	if v.UnselectedTextStyle != nil {
		unselectedTextStyle = v.UnselectedTextStyle.MarshalProtobuf()
	}

	return view.Model{
		Children:       l.Views(),
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/tabscreen",
		NativeViewState: internal.MarshalProtobuf(&pbios.TabView{
			Screens:             childrenPb,
			SelectedIndex:       int64(v.Tabs.SelectedIndex()),
			BarColor:            pb.ColorEncode(v.BarColor),
			SelectedColor:       pb.ColorEncode(v.SelectedColor),
			UnselectedColor:     pb.ColorEncode(v.UnselectedColor),
			SelectedTextStyle:   selectedTextStyle,
			UnselectedTextStyle: unselectedTextStyle,
		}),
		NativeFuncs: map[string]interface{}{
			"OnSelect": func(data []byte) {
				pbevent := &pbios.TabEvent{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Tabs.SetSelectedIndex(int(pbevent.SelectedIndex))
			},
		},
	}
}

// TabButton describes a UITabBarItem.
type TabButton struct {
	Title        string
	Icon         image.Image
	SelectedIcon image.Image
	Badge        string
}

func (t *TabButton) OptionKey() string {
	return "gomatcha.io/view/ios TabButton"
}

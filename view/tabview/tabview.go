package tabview

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/pb"
	pbtext "gomatcha.io/matcha/pb/text"
	tabnavpb "gomatcha.io/matcha/pb/view/tabscreen"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

type Tabs struct {
	relay         comm.Relay
	children      []view.View
	selectedIndex int
}

func (s *Tabs) SetViews(ss ...view.View) {
	s.children = ss
	s.relay.Signal()
}

func (s *Tabs) Views() []view.View {
	return s.children
}

func (s *Tabs) SetSelectedIndex(idx int) {
	if idx != s.selectedIndex {
		s.selectedIndex = idx
		s.relay.Signal()
	}
}

func (s *Tabs) SelectedIndex() int {
	return s.selectedIndex
}

func (s *Tabs) Notify(f func()) comm.Id {
	return s.relay.Notify(f)
}

func (s *Tabs) Unnotify(id comm.Id) {
	s.relay.Unnotify(id)
}

type View struct {
	view.Embed
	Tabs                *Tabs
	BarColor            color.Color
	SelectedTextStyle   *text.Style
	UnselectedTextStyle *text.Style
	SelectedColor       color.Color
	UnselectedColor     color.Color
	tabs                *Tabs
}

func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *View) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		if v.tabs != nil {
			v.Unsubscribe(v.tabs)
		}
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	// Subscribe to the group
	if v.Tabs != v.tabs {
		if v.tabs != nil {
			v.Unsubscribe(v.tabs)
		}
		if v.Tabs != nil {
			v.Subscribe(v.Tabs)
		}
		v.tabs = v.Tabs
	}

	childrenPb := []*tabnavpb.ChildView{}
	for _, chld := range v.Tabs.Views() {
		// Create the button
		var button *Button
		if childView, ok := chld.(ChildView); ok {
			button = childView.TabButton(ctx)
		} else {
			button = &Button{
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
		childrenPb = append(childrenPb, &tabnavpb.ChildView{
			Id:           int64(chld.Id()),
			Title:        button.Title,
			Icon:         app.ImageMarshalProtobuf(button.Icon),
			SelectedIcon: app.ImageMarshalProtobuf(button.SelectedIcon),
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
		NativeViewState: &tabnavpb.View{
			Screens:             childrenPb,
			SelectedIndex:       int64(v.Tabs.SelectedIndex()),
			BarColor:            pb.ColorEncode(v.BarColor),
			SelectedColor:       pb.ColorEncode(v.SelectedColor),
			UnselectedColor:     pb.ColorEncode(v.UnselectedColor),
			SelectedTextStyle:   selectedTextStyle,
			UnselectedTextStyle: unselectedTextStyle,
		},
		NativeFuncs: map[string]interface{}{
			"OnSelect": func(data []byte) {
				pbevent := &tabnavpb.Event{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Tabs.SetSelectedIndex(int(pbevent.SelectedIndex))
			},
		},
	}

	// children := map[int64]view.View{}
	// childrenPb := []*tabnavpb.ChildView{}
	// v.ids = append([]int64(nil), v.Group.ids...)
	// for _, i := range v.ids {
	// 	key := strconv.Itoa(int(i))

	// 	// Create the child if necessary and subscribe to it.
	// 	chld, ok := v.children[i]
	// 	if !ok {
	// 		chld = v.Group.children[i].View(ctx.WithPrefix("view" + key))
	// 		children[i] = chld
	// 		v.Subscribe(chld)
	// 	} else {
	// 		children[i] = chld
	// 		delete(v.children, i)
	// 	}

	// // Create the button
	// var button *Button
	// if childView, ok := chld.(ChildView); ok {
	// 	button = childView.TabButton(ctx)
	// } else {
	// 	button = &Button{
	// 		Title: "Title",
	// 	}
	// }

	// // Add the child.
	// l.Add(chld, func(s *constraint.Solver) {
	// 	s.TopEqual(constraint.Const(0))
	// 	s.LeftEqual(constraint.Const(0))
	// 	s.WidthEqual(l.MaxGuide().Width())
	// 	s.HeightEqual(l.MaxGuide().Height())
	// })

	// // Add to protobuf.
	// childrenPb = append(childrenPb, &tabnavpb.ChildView{
	// 	Id:           int64(chld.Id()),
	// 	Title:        button.Title,
	// 	Icon:         env.ImageMarshalProtobuf(button.Icon),
	// 	SelectedIcon: env.ImageMarshalProtobuf(button.SelectedIcon),
	// 	Badge:        button.Badge,
	// })
	// }

	// // Unsubscribe from old views
	// for _, chld := range v.children {
	// 	v.Unsubscribe(chld)
	// }
	// v.children = children
}

type ChildView interface {
	view.View
	TabButton(*view.Context) *Button
}

type Button struct {
	Title        string
	Icon         image.Image
	SelectedIcon image.Image
	Badge        string
}

func WithButton(s view.View, button *Button) view.View {
	return &viewWrapper{
		View:   s,
		button: button,
	}
}

type viewWrapper struct {
	view.View
	button *Button
}

func (v *viewWrapper) TabButton(*view.Context) *Button {
	return v.button
}

package stackview

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/pb"
	pbtext "gomatcha.io/matcha/pb/text"
	"gomatcha.io/matcha/pb/view/stacknav"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

type Stack struct {
	relay    comm.Relay
	children []view.View
}

func (s *Stack) SetViews(ss ...view.View) {
	s.children = ss
	s.relay.Signal()
}

func (s *Stack) setChildIds(ids []int64) {
	prevChildren := s.children
	s.children = []view.View{}

	for _, i := range ids {
		for _, j := range prevChildren {
			if j.Id() == matcha.Id(i) {
				s.children = append(s.children, j)
				continue
			}
		}
	}
	s.relay.Signal()
}

func (s *Stack) Views() []view.View {
	return s.children
}

func (s *Stack) Push(vs view.View) {
	s.children = append(s.children, vs)
	s.relay.Signal()
}

func (s *Stack) Pop() {
	s.children = s.children[:len(s.children)-1]
	s.relay.Signal()
}

func (s *Stack) Notify(f func()) comm.Id {
	return s.relay.Notify(f)
}

func (s *Stack) Unnotify(id comm.Id) {
	s.relay.Unnotify(id)
}

type View struct {
	view.Embed
	Stack          *Stack
	stack          *Stack
	TitleTextStyle *text.Style
	BackTextStyle  *text.Style
	BarColor       color.Color
	// children map[int64]view.View
	// ids      []int64
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
		if v.stack != nil {
			v.Unsubscribe(v.stack)
		}
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	// Subscribe to the stack
	if v.Stack != v.stack {
		if v.stack != nil {
			v.Unsubscribe(v.stack)
		}
		if v.Stack != nil {
			v.Subscribe(v.Stack)
		}
		v.stack = v.Stack
	}

	childrenPb := []*stacknav.ChildView{}
	for _, chld := range v.Stack.Views() {
		key := strconv.Itoa(int(chld.Id()))

		// v.Subscribe(chld)

		// Create the bar.
		var bar *Bar
		if childView, ok := chld.(ChildView); ok {
			bar = childView.StackBar(ctx.WithPrefix("bar" + key))
		} else {
			bar = &Bar{
				Title: "Title",
			}
		}

		// Add the bar.
		barV := &barView{
			Embed: ctx.NewEmbed(key),
			bar:   bar,
		}
		l.Add(barV, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.WidthEqual(l.MaxGuide().Width())
			s.HeightEqual(constraint.Const(44))
		})

		// Add the child.
		l.Add(chld, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.WidthEqual(l.MaxGuide().Width())
			s.HeightEqual(l.MaxGuide().Height().Add(-64)) // TODO(KD): Respect bar actual height, shorter when rotated, etc...
		})

		// Add ids to protobuf.
		childrenPb = append(childrenPb, &stacknav.ChildView{
			ViewId:   int64(chld.Id()),
			BarId:    int64(barV.Id()),
			ScreenId: int64(chld.Id()),
		})
	}

	var titleTextStyle *pbtext.TextStyle
	if v.TitleTextStyle != nil {
		titleTextStyle = v.TitleTextStyle.MarshalProtobuf()
	}

	var backTextStyle *pbtext.TextStyle
	if v.BackTextStyle != nil {
		backTextStyle = v.BackTextStyle.MarshalProtobuf()
	}

	return view.Model{
		Children:       l.Views(),
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/stacknav",
		NativeViewState: &stacknav.View{
			Children:       childrenPb,
			TitleTextStyle: titleTextStyle,
			BackTextStyle:  backTextStyle,
			BarColor:       pb.ColorEncode(v.BarColor),
		},
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				pbevent := &stacknav.StackEvent{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Stack.setChildIds(pbevent.Id)
			},
		},
	}

	// children := map[int64]view.View{}
	// childrenPb := []*stacknav.ChildView{}
	// v.ids = append([]int64(nil), v.stack.ids...)
	// for _, i := range v.ids {
	// 	key := strconv.Itoa(int(i))

	// 	// Create the child if necessary and subscribe to it.
	// 	chld, ok := v.children[i]
	// 	if !ok {
	// 		chld = v.stack.children[i].View(ctx.WithPrefix("view" + key))
	// 		children[i] = chld
	// 		v.Subscribe(chld)
	// 	} else {
	// 		children[i] = chld
	// 		delete(v.children, i)
	// 	}

	// // Create the bar.
	// var bar *Bar
	// if childView, ok := chld.(ChildView); ok {
	// 	bar = childView.StackBar(ctx.WithPrefix("bar" + key))
	// } else {
	// 	bar = &Bar{
	// 		Title: "Title",
	// 	}
	// }

	// // Add the bar.
	// barV := &barView{
	// 	Embed: view.NewEmbed(ctx.NewId(key)),
	// 	bar:   bar,
	// }
	// l.Add(barV, func(s *constraint.Solver) {
	// 	s.TopEqual(constraint.Const(0))
	// 	s.LeftEqual(constraint.Const(0))
	// 	s.WidthEqual(l.MaxGuide().Width())
	// 	s.HeightEqual(constraint.Const(44))
	// })

	// // Add the child.
	// l.Add(chld, func(s *constraint.Solver) {
	// 	s.TopEqual(constraint.Const(0))
	// 	s.LeftEqual(constraint.Const(0))
	// 	s.WidthEqual(l.MaxGuide().Width())
	// 	s.HeightEqual(l.MaxGuide().Height().Add(-64)) // TODO(KD): Respect bar actual height, shorter when rotated, etc...
	// })

	// // Add ids to protobuf.
	// childrenPb = append(childrenPb, &stacknav.ChildView{
	// 	ViewId:   int64(chld.Id()),
	// 	BarId:    int64(barV.Id()),
	// 	ScreenId: i,
	// })
	// }

	// // Unsubscribe from old views
	// for _, chld := range v.children {
	// 	v.Unsubscribe(chld)
	// }
	// v.children = children

	// return view.Model{
	// 	Children:       l.Views(),
	// 	Layouter:       l,
	// 	NativeViewName: "gomatcha.io/matcha/view/stacknav",
	// 	NativeViewState: &stacknav.View{
	// 		Children: childrenPb,
	// 	},
	// 	NativeFuncs: map[string]interface{}{
	// 		"OnChange": func(data []byte) {
	// 			pbevent := &stacknav.StackEvent{}
	// 			err := proto.Unmarshal(data, pbevent)
	// 			if err != nil {
	// 				fmt.Println("error", err)
	// 				return
	// 			}

	// 			v.stack.Lock()
	// 			v.stack.setChildIds(pbevent.Id)
	// 			v.stack.Unlock()
	// 		},
	// 	},
	// }
}

type ChildView interface {
	view.View
	StackBar(*view.Context) *Bar // TODO(KD): Doesn't this make it harder to wrap??
}

type barView struct {
	view.Embed
	bar *Bar
}

func (v *barView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	// iOS does the layouting for us. We just need the correct sizes.
	titleViewId := int64(0)
	if v.bar.TitleView != nil {
		titleViewId = int64(v.bar.TitleView.Id())
		l.Add(v.bar.TitleView, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.HeightLess(l.MaxGuide().Height())
			s.WidthLess(l.MaxGuide().Width())
		})
	}

	rightViewIds := []int64{}
	for _, i := range v.bar.RightViews {
		rightViewIds = append(rightViewIds, int64(i.Id()))
		l.Add(i, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.HeightLess(l.MaxGuide().Height())
			s.WidthLess(l.MaxGuide().Width())
		})
	}
	leftViewIds := []int64{}
	for _, i := range v.bar.LeftViews {
		leftViewIds = append(leftViewIds, int64(i.Id()))
		l.Add(i, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(0))
			s.LeftEqual(constraint.Const(0))
			s.HeightLess(l.MaxGuide().Height())
			s.WidthLess(l.MaxGuide().Width())
		})
	}

	return view.Model{
		Layouter:       l,
		Children:       l.Views(),
		NativeViewName: "gomatcha.io/matcha/view/stacknav Bar",
		NativeViewState: &stacknav.Bar{
			Title: v.bar.Title,
			CustomBackButtonTitle: len(v.bar.BackButtonTitle) > 0,
			BackButtonTitle:       v.bar.BackButtonTitle,
			BackButtonHidden:      v.bar.BackButtonHidden,
			TitleViewId:           titleViewId,
			RightViewIds:          rightViewIds,
			LeftViewIds:           leftViewIds,
		},
	}
}

type Bar struct {
	Title            string
	BackButtonTitle  string
	BackButtonHidden bool

	TitleView  view.View
	RightViews []view.View
	LeftViews  []view.View
}

func WithBar(s view.View, bar *Bar) view.View {
	return &viewWrapper{
		View:     s,
		stackBar: bar,
	}
}

type viewWrapper struct {
	view.View
	stackBar *Bar
}

func (s *viewWrapper) StackBar(*view.Context) *Bar {
	return s.stackBar
}

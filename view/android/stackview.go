package android

import (
	"strconv"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/pb/view/android"
	"gomatcha.io/matcha/view"
)

// Stack represents a list of views to be shown in the StackView. It can be manipulated outside of a Build() call.
type Stack struct {
	relay       comm.Relay
	childIds    []int64
	childrenMap map[int64]view.View
	maxId       int64
}

func (s *Stack) SetViews(vs ...view.View) {
	if s.childrenMap == nil {
		s.childrenMap = map[int64]view.View{}
	}

	for _, i := range vs {
		s.maxId += 1
		s.childIds = append(s.childIds, s.maxId)
		s.childrenMap[s.maxId] = i
	}
	s.relay.Signal()
}

func (s *Stack) setChildIds(ids []int64) {
	s.childIds = ids
	s.relay.Signal()
}

func (s *Stack) Views() []view.View {
	vs := []view.View{}
	for _, i := range s.childIds {
		vs = append(vs, s.childrenMap[i])
	}
	return vs
}

func (s *Stack) Push(vs view.View) {
	s.maxId += 1

	s.childIds = append(s.childIds, s.maxId)
	s.childrenMap[s.maxId] = vs
	s.relay.Signal()
}

func (s *Stack) Pop() {
	delete(s.childrenMap, s.childIds[len(s.childIds)-1])
	s.childIds = s.childIds[:len(s.childIds)-1]
	s.relay.Signal()
}

func (s *Stack) Notify(f func()) comm.Id {
	return s.relay.Notify(f)
}

func (s *Stack) Unnotify(id comm.Id) {
	s.relay.Unnotify(id)
}

type StackView struct {
	view.Embed
	Stack *Stack
	stack *Stack
	// TitleTextStyle *text.Style
	// BackTextStyle  *text.Style
	// BarColor       color.Color
	// children map[int64]view.View
	// ids      []int64
}

// NewStackView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewStackView() *StackView {
	return &StackView{
		Stack: &Stack{},
	}
}

// Lifecyle implements the view.View interface.
func (v *StackView) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		if v.stack != nil {
			v.Unsubscribe(v.stack)
		}
	}
}

// Build implements the view.View interface.
func (v *StackView) Build(ctx *view.Context) view.Model {
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

	childrenPb := []*android.StackChildView{}
	for _, id := range v.Stack.childIds {
		chld := v.Stack.childrenMap[id]

		// Find the bar.
		var bar *StackBar
		for _, opts := range chld.Build(nil).Options {
			var ok bool
			if bar, ok = opts.(*StackBar); ok {
				break
			}
		}
		if bar == nil {
			bar = &StackBar{
				Title: "Title",
			}
		}

		// Add the bar.
		barV := &stackBarView{
			Embed: view.Embed{Key: strconv.Itoa(int(id))},
			Bar:   bar,
		}
		l.Add(barV, func(s *constraint.Solver) {
			s.Top(0)
			s.Left(0)
			s.WidthEqual(l.MaxGuide().Width())
			s.Height(44)
		})

		// Add the child.
		l.Add(chld, func(s *constraint.Solver) {
			s.Top(0)
			s.Left(0)
			s.WidthEqual(l.MaxGuide().Width())
			s.HeightEqual(l.MaxGuide().Height().Add(-64)) // TODO(KD): Respect bar actual height, shorter when rotated, etc...
		})

		// Add ids to protobuf.
		childrenPb = append(childrenPb, &android.StackChildView{
			ScreenId: int64(id),
		})
	}

	return view.Model{
		Children:       l.Views(),
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/android StackView",
		NativeViewState: &android.StackView{
			Children: childrenPb,
			// BarColor: pb.ColorEncode(v.BarColor),
		},
		NativeFuncs: map[string]interface{}{
		// "OnChange": func(data []byte) {
		// 	pbevent := &stacknav.StackEvent{}
		// 	err := proto.Unmarshal(data, pbevent)
		// 	if err != nil {
		// 		fmt.Println("error", err)
		// 		return
		// 	}

		// 	v.Stack.setChildIds(pbevent.Id)
		// },
		},
	}
}

type stackBarView struct {
	view.Embed
	Bar *StackBar
}

func (v *stackBarView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	// // iOS does the layouting for us. We just need the correct sizes.
	// hasTitleView := false
	// if v.Bar.TitleView != nil {
	// 	hasTitleView = true
	// 	l.Add(v.Bar.TitleView, func(s *constraint.Solver) {
	// 		s.Top(0)
	// 		s.Left(0)
	// 		s.HeightLess(l.MaxGuide().Height())
	// 		s.WidthLess(l.MaxGuide().Width())
	// 	})
	// }

	// rightViewCount := int64(0)
	// for _, i := range v.Bar.RightViews {
	// 	rightViewCount += 1
	// 	l.Add(i, func(s *constraint.Solver) {
	// 		s.Top(0)
	// 		s.Left(0)
	// 		s.HeightLess(l.MaxGuide().Height())
	// 		s.WidthLess(l.MaxGuide().Width())
	// 	})
	// }
	// leftViewCount := int64(0)
	// for _, i := range v.Bar.LeftViews {
	// 	leftViewCount += 1
	// 	l.Add(i, func(s *constraint.Solver) {
	// 		s.Top(0)
	// 		s.Left(0)
	// 		s.HeightLess(l.MaxGuide().Height())
	// 		s.WidthLess(l.MaxGuide().Width())
	// 	})
	// }

	return view.Model{
		Layouter:       l,
		Children:       l.Views(),
		NativeViewName: "gomatcha.io/matcha/view/android stackBarView",
		NativeViewState: &android.StackBar{
			Title: v.Bar.Title,
			// CustomBackButtonTitle: len(v.Bar.BackButtonTitle) > 0,
			// BackButtonTitle:       v.Bar.BackButtonTitle,
			// BackButtonHidden:      v.Bar.BackButtonHidden,
			// HasTitleView:          hasTitleView,
			// RightViewCount:        rightViewCount,
			// LeftViewCount:         leftViewCount,
		},
	}
}

type StackBar struct {
	Title string
	// BackButtonTitle  string
	// BackButtonHidden bool

	TitleView  view.View
	RightViews []view.View
	LeftViews  []view.View
}

func (t *StackBar) OptionKey() string {
	return "gomatcha.io/view/android StackBar"
}

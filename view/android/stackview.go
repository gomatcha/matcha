package android

import (
	"image"
	"image/color"
	"strconv"

	"github.com/gogo/protobuf/proto"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	"gomatcha.io/matcha/pb/view/android"
	"gomatcha.io/matcha/text"
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
	if len(s.childIds) <= 1 {
		return
	}
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
func (v *StackView) Build(ctx view.Context) view.Model {
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
	for idx, id := range v.Stack.childIds {
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
			Embed:           view.Embed{Key: strconv.Itoa(int(id))},
			Bar:             bar,
			NeedsBackButton: idx != 0,
		}
		l.Add(barV, func(s *constraint.Solver) {
			s.Top(0)
			s.Left(0)
			s.WidthEqual(l.MaxGuide().Width())
			s.Height(56)
		})

		// Add the child.
		l.Add(chld, func(s *constraint.Solver) {
			s.Top(0)
			s.Left(0)
			s.WidthEqual(l.MaxGuide().Width())
			s.HeightEqual(l.MaxGuide().Height().Add(-56 - 24)) // TODO(KD): Respect bar actual height, shorter when rotated, etc...
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
			"OnBack": func() {
				v.Stack.Pop()
			},
			"CanBack": func() bool {
				return len(v.Stack.childIds) >= 2
			},
		},
	}
}

type stackBarView struct {
	view.Embed
	Bar             *StackBar
	NeedsBackButton bool
}

func (v *stackBarView) Build(ctx view.Context) view.Model {
	funcs := map[string]interface{}{}
	items := []*android.StackBarItem{}
	for idx, i := range v.Bar.Buttons {
		button := i.marshalProtobuf()
		button.OnPressFunc = strconv.Itoa(idx)
		items = append(items, button)
		funcs[strconv.Itoa(idx)] = i.OnPress
	}

	return view.Model{
		Painter:        &paint.Style{BackgroundColor: v.Bar.Color},
		NativeViewName: "gomatcha.io/matcha/view/android stackBarView",
		NativeViewState: &android.StackBar{
			Title:            v.Bar.Title,
			StyledTitle:      v.Bar.StyledTitle.MarshalProtobuf(),
			Subtitle:         v.Bar.Subtitle,
			StyledSubtitle:   v.Bar.StyledSubtitle.MarshalProtobuf(),
			Items:            items,
			BackButtonHidden: !v.NeedsBackButton,
		},
		NativeFuncs: funcs,
	}
}

type StackBar struct {
	Title          string
	StyledTitle    *text.StyledText
	Subtitle       string
	StyledSubtitle *text.StyledText
	Color          color.Color
	Buttons        []*StackBarButton
}

func (t *StackBar) OptionKey() string {
	return "gomatcha.io/view/android StackBar"
}

type StackBarButton struct {
	Title       string
	StyledTitle *text.StyledText
	Icon        image.Image
	IconTint    color.Color
	Disabled    bool
	OnPress     func()
}

func (v *StackBarButton) marshalProtobuf() *android.StackBarItem {
	return &android.StackBarItem{
		Title:    v.Title,
		Icon:     app.ImageMarshalProtobuf(v.Icon),
		IconTint: pb.ColorEncode(v.IconTint),
		Disabled: v.Disabled,
	}
}

var stackMiddlewareVar *stackMiddleware

func init() {
	internal.RegisterMiddleware(func() interface{} {
		stackMiddlewareVar = &stackMiddleware{
			radix: radix.NewRadix(),
		}
		return stackMiddlewareVar
	})
	bridge.RegisterFunc("gomatcha.io/view/android StackBarOnBack", func() {
		didBack := false
		stackMiddlewareVar.radix.Range(func(path []int64, node *radix.Node) {
			if !didBack {
				didBack = true
				node.Value.(map[string]interface{})["OnBack"].(func())()
			}
		})
	})
	bridge.RegisterFunc("gomatcha.io/view/android StackBarCanBack", func() bool {
		canBack := false
		stackMiddlewareVar.radix.Range(func(path []int64, node *radix.Node) {
			canBack = node.Value.(map[string]interface{})["CanBack"].(func() bool)()
		})
		return canBack
	})
}

type stackMiddleware struct {
	radix *radix.Radix
}

func (m *stackMiddleware) Build(ctx view.Context, model *view.Model) {
	path := idSliceToIntSlice(ctx.Path())

	var nativeFuncs map[string]interface{}
	if model.NativeViewName == "gomatcha.io/matcha/view/android StackView" {
		nativeFuncs = model.NativeFuncs
	}

	if nativeFuncs != nil {
		n := m.radix.Insert(path)
		n.Value = nativeFuncs
	} else {
		m.radix.Delete(path)
	}
}

func (m *stackMiddleware) MarshalProtobuf() proto.Message {
	return nil
}

func (m *stackMiddleware) Key() string {
	return "gomatcha.io/matcha/view/android stackMiddleware"
}

func idSliceToIntSlice(ids []view.Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

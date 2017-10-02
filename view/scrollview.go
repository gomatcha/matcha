package view

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"

	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pbview "gomatcha.io/matcha/proto/view"
)

type ScrollView struct {
	Embed
	ScrollAxes     layout.Axis
	IndicatorAxes  layout.Axis
	ScrollEnabled  bool
	ScrollPosition *ScrollPosition
	OnScroll       func(position layout.Point)

	ContentChildren []View
	ContentPainter  paint.Painter
	ContentLayouter layout.Layouter
	PaintStyle      *paint.Style
}

// NewScrollView returns a new view.
func NewScrollView() *ScrollView {
	return &ScrollView{
		ScrollPosition: &ScrollPosition{},
		ScrollAxes:     layout.AxisY,
		IndicatorAxes:  layout.AxisY | layout.AxisX,
		ScrollEnabled:  true,
	}
}

func (v *ScrollView) Lifecycle(from, to Stage) {
	if EntersStage(from, to, StageMounted) {
		if v.ScrollPosition == nil {
			v.ScrollPosition = &ScrollPosition{}
		}
	}
}

func (v *ScrollView) Update(v2 View) {
	CopyFields(v, v2)

	if v.ScrollPosition == nil {
		v.ScrollPosition = &ScrollPosition{}
	}
}

// Build implements View.
func (v *ScrollView) Build(ctx Context) Model {
	child := NewBasicView()
	child.Children = v.ContentChildren
	child.Layouter = v.ContentLayouter
	child.Painter = v.ContentPainter

	var painter paint.Painter
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Children: []View{child},
		Painter:  painter,
		Layouter: &scrollViewLayouter{
			axes:           v.ScrollAxes,
			scrollPosition: v.ScrollPosition,
		},
		NativeViewName: "gomatcha.io/matcha/view/scrollview",
		NativeViewState: internal.MarshalProtobuf(&pbview.ScrollView{
			ScrollEnabled:                  v.ScrollEnabled,
			Horizontal:                     v.ScrollAxes|layout.AxisX == layout.AxisX,
			Vertical:                       v.ScrollAxes|layout.AxisY == layout.AxisY,
			ShowsHorizontalScrollIndicator: v.IndicatorAxes&layout.AxisY == layout.AxisY,
			ShowsVerticalScrollIndicator:   v.IndicatorAxes&layout.AxisX == layout.AxisX,
		}),
		NativeFuncs: map[string]interface{}{
			"OnScroll": func(data []byte) {
				event := &pbview.ScrollEvent{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				var offset layout.Point
				(&offset).UnmarshalProtobuf(event.ContentOffset)

				v.ScrollPosition.setValue(offset, true)

				if v.OnScroll != nil {
					v.OnScroll(offset)
				}
			},
		},
	}
}

type scrollViewLayouter struct {
	axes           layout.Axis
	scrollPosition *ScrollPosition
}

func (l *scrollViewLayouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	minSize := ctx.MinSize()
	maxSize := ctx.MaxSize()
	if l.axes&layout.AxisY == layout.AxisY {
		minSize.Y = 0
		maxSize.Y = math.Inf(1)
	}
	if l.axes&layout.AxisX == layout.AxisX {
		minSize.X = 0
		maxSize.X = math.Inf(1)
	}

	offset := l.scrollPosition.Value()
	g := ctx.LayoutChild(0, minSize, maxSize)
	g.Frame = layout.Rt(-offset.X, -offset.Y, g.Width()-offset.X, g.Height()-offset.Y)
	gs := []layout.Guide{g}

	return layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize().X, ctx.MinSize().Y),
	}, gs
}

func (l *scrollViewLayouter) Notify(f func()) comm.Id {
	return l.scrollPosition.notify(f, false)
}

func (l *scrollViewLayouter) Unnotify(id comm.Id) {
	l.scrollPosition.Unnotify(id)
}

type ScrollPosition struct {
	X           animate.Value
	Y           animate.Value
	group       comm.Relay
	initialized bool
	userEvent   bool
}

func (p *ScrollPosition) initialize() {
	if p.initialized {
		return
	}
	p.initialized = true
	p.group.Subscribe(&p.X)
	p.group.Subscribe(&p.Y)
}

func (p *ScrollPosition) Notify(f func()) comm.Id {
	return p.notify(f, true)
}

// Allow ignoring user triggered scroll events.
func (p *ScrollPosition) notify(f func(), userEvents bool) comm.Id {
	p.initialize()
	if userEvents {
		return p.group.Notify(f)
	} else {
		return p.group.Notify(func() {
			if !p.userEvent {
				f()
			}
		})
	}
}

func (p *ScrollPosition) Unnotify(id comm.Id) {
	p.initialize()
	p.group.Unnotify(id)
}

func (p *ScrollPosition) Value() layout.Point {
	return layout.Pt(p.X.Value(), p.Y.Value())
}

func (p *ScrollPosition) SetValue(val layout.Point) {
	p.setValue(val, false)
}

func (p *ScrollPosition) setValue(val layout.Point, userEvent bool) {
	p.userEvent = userEvent
	if val != p.Value() {
		p.X.SetValue(val.X)
		p.Y.SetValue(val.Y)
	}
	p.userEvent = false
}

// TODO(KD):
// func (p *ScrollPosition) ScrollToPoint(val layout.Point) {
// }

// func (p *ScrollPosition) ScrollToChild(id comm.Id) {
// }

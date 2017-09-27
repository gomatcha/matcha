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
	scrollPosition *ScrollPosition
	offset         *layout.Point
	OnScroll       func(position layout.Point)

	ContentChildren []View
	ContentPainter  paint.Painter
	ContentLayouter layout.Layouter
	PaintStyle      *paint.Style
}

// NewScrollView returns a new view.
func NewScrollView() *ScrollView {
	return &ScrollView{
		ScrollAxes:    layout.AxisY,
		IndicatorAxes: layout.AxisY | layout.AxisX,
		ScrollEnabled: true,
		offset:        &layout.Point{},
	}
}

// Build implements view.View.
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
			offset:         v.offset,
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

				*v.offset = offset
				if v.ScrollPosition != nil {
					v.ScrollPosition.SetValue(offset)
				}
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
	offset         *layout.Point
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

	g := ctx.LayoutChild(0, minSize, maxSize)
	g.Frame = layout.Rt(-l.offset.X, -l.offset.Y, g.Width()-l.offset.X, g.Height()-l.offset.Y)
	gs := []layout.Guide{g}

	return layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize().X, ctx.MinSize().Y),
	}, gs
}

func (l *scrollViewLayouter) Notify(f func()) comm.Id {
	if l.scrollPosition == nil {
		return 0
	}
	return l.scrollPosition.Notify(func() {
		if *l.offset != l.scrollPosition.Value() {
			*l.offset = l.scrollPosition.Value()
			f()
		}
	})
}

func (l *scrollViewLayouter) Unnotify(id comm.Id) {
	if l.scrollPosition == nil {
		return
	}
	l.scrollPosition.Unnotify(id)
}

type ScrollPosition struct {
	X           animate.Value
	Y           animate.Value
	group       comm.Relay
	initialized bool
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
	p.initialize()
	return p.group.Notify(f)
}

func (p *ScrollPosition) Unnotify(id comm.Id) {
	p.initialize()
	p.group.Unnotify(id)
}

func (p *ScrollPosition) Value() layout.Point {
	return layout.Pt(p.X.Value(), p.Y.Value())
}

func (p *ScrollPosition) SetValue(val layout.Point) {
	if val == p.Value() {
		return
	}
	p.X.SetValue(val.X)
	p.Y.SetValue(val.Y)
}

// TODO(KD):
// func (p *ScrollPosition) ScrollToPoint(val layout.Point) {
// }

// func (p *ScrollPosition) ScrollToChild(id comm.Id) {
// }

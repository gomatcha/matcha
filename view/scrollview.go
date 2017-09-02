package view

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"

	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/scrollview"
)

// ScrollDirection represents the X and Y axis.
type ScrollDirection int

const (
	ScrollDirectionHorizontal ScrollDirection = 1 << iota
	ScrollDirectionVertical
)

type ScrollView struct {
	Embed
	Directions          ScrollDirection
	IndicatorDirections ScrollDirection
	ScrollEnabled       bool
	ScrollPosition      *ScrollPosition
	scrollPosition      *ScrollPosition
	offset              *layout.Point
	OnScroll            func(position layout.Point)

	ContentChildren []View
	ContentPainter  paint.Painter
	ContentLayouter layout.Layouter
	PaintStyle      *paint.Style
}

// NewScrollView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewScrollView() *ScrollView {
	return &ScrollView{
		Directions:          ScrollDirectionVertical,
		IndicatorDirections: ScrollDirectionVertical | ScrollDirectionHorizontal,
		ScrollEnabled:       true,
		offset:              &layout.Point{},
	}
}

// Build implements view.View.
func (v *ScrollView) Build(ctx *Context) Model {
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
			directions:     v.Directions,
			scrollPosition: v.ScrollPosition,
			offset:         v.offset,
		},
		NativeViewName: "gomatcha.io/matcha/view/scrollview",
		NativeViewState: &scrollview.View{
			ScrollEnabled:                  v.ScrollEnabled,
			ShowsHorizontalScrollIndicator: v.IndicatorDirections&ScrollDirectionHorizontal == ScrollDirectionHorizontal,
			ShowsVerticalScrollIndicator:   v.IndicatorDirections&ScrollDirectionVertical == ScrollDirectionVertical,
		},
		NativeFuncs: map[string]interface{}{
			"OnScroll": func(data []byte) {
				event := &scrollview.ScrollEvent{}
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
	directions     ScrollDirection
	scrollPosition *ScrollPosition
	offset         *layout.Point
}

func (l *scrollViewLayouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	minSize := ctx.MinSize
	if l.directions&ScrollDirectionHorizontal == ScrollDirectionHorizontal {
		minSize.X = 0
	}
	if l.directions&ScrollDirectionVertical == ScrollDirectionVertical {
		minSize.Y = 0
	}

	g := ctx.LayoutChild(0, minSize, layout.Pt(math.Inf(1), math.Inf(1)))
	g.Frame = layout.Rt(-l.offset.X, -l.offset.Y, g.Width()-l.offset.X, g.Height()-l.offset.Y)
	gs := []layout.Guide{g}

	return layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
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

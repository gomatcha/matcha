// Package button implements a native scroll view.
package scrollview

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"

	"gomatcha.io/matcha"
	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/scrollview"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
)

// Direction represents the X and Y axis.
type Direction int

const (
	Horizontal Direction = 1 << iota
	Vertical
)

type View struct {
	view.Embed
	Direction                Direction
	ScrollIndicatorDirection Direction
	ScrollEnabled            bool
	ScrollPosition           *ScrollPosition
	scrollPosition           *ScrollPosition
	OnScroll                 func(position layout.Point)

	ContentChildren []view.View
	ContentPainter  paint.Painter
	ContentLayouter layout.Layouter
	PaintStyle      *paint.Style
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed:                    ctx.NewEmbed(key),
		Direction:                Vertical,
		ScrollIndicatorDirection: Vertical | Horizontal,
		ScrollEnabled:            true,
		scrollPosition:           &ScrollPosition{},
	}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	child := basicview.New(ctx, "child")
	child.Children = v.ContentChildren
	child.Layouter = v.ContentLayouter
	child.Painter = v.ContentPainter

	scrollPosition := v.ScrollPosition
	if scrollPosition == nil {
		scrollPosition = v.scrollPosition
	}

	var painter paint.Painter
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Children: []view.View{child},
		Painter:  painter,
		Layouter: &layouter{
			directions:     v.Direction,
			scrollPosition: scrollPosition,
		},
		NativeViewName: "gomatcha.io/matcha/view/scrollview",
		NativeViewState: &scrollview.View{
			ScrollEnabled:                  v.ScrollEnabled,
			ShowsHorizontalScrollIndicator: v.ScrollIndicatorDirection&Horizontal == Horizontal,
			ShowsVerticalScrollIndicator:   v.ScrollIndicatorDirection&Vertical == Vertical,
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

				v.scrollPosition.SetValue(offset)
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

type layouter struct {
	directions     Direction
	scrollPosition *ScrollPosition // TODO(KD): Are we unnotifying this correctly?
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	gs := map[matcha.Id]layout.Guide{}

	minSize := ctx.MinSize
	if l.directions&Horizontal == Horizontal {
		minSize.X = 0
	}
	if l.directions&Vertical == Vertical {
		minSize.Y = 0
	}

	position := l.scrollPosition.Value()
	g := ctx.LayoutChild(ctx.ChildIds[0], minSize, layout.Pt(math.Inf(1), math.Inf(1)))
	g.Frame = layout.Rt(-position.X, -position.Y, g.Width()-position.X, g.Height()-position.Y)
	gs[ctx.ChildIds[0]] = g

	return layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
	}, gs
}

func (l *layouter) Notify(f func()) comm.Id {
	return l.scrollPosition.Notify(f)
}

func (l *layouter) Unnotify(id comm.Id) {
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

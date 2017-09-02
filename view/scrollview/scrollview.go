// Package scrollview implements a native scroll view.
package scrollview

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"

	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/scrollview"
	"gomatcha.io/matcha/view"
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
	offset                   *layout.Point
	OnScroll                 func(position layout.Point)

	ContentChildren []view.View
	ContentPainter  paint.Painter
	ContentLayouter layout.Layouter
	PaintStyle      *paint.Style
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{
		Direction:                Vertical,
		ScrollIndicatorDirection: Vertical | Horizontal,
		ScrollEnabled:            true,
		offset:                   &layout.Point{},
	}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	child := view.NewBasicView()
	child.Children = v.ContentChildren
	child.Layouter = v.ContentLayouter
	child.Painter = v.ContentPainter

	var painter paint.Painter
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Children: []view.View{child},
		Painter:  painter,
		Layouter: &layouter{
			directions:     v.Direction,
			scrollPosition: v.ScrollPosition,
			offset:         v.offset,
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

type layouter struct {
	directions     Direction
	scrollPosition *ScrollPosition
	offset         *layout.Point
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	minSize := ctx.MinSize
	if l.directions&Horizontal == Horizontal {
		minSize.X = 0
	}
	if l.directions&Vertical == Vertical {
		minSize.Y = 0
	}

	g := ctx.LayoutChild(0, minSize, layout.Pt(math.Inf(1), math.Inf(1)))
	g.Frame = layout.Rt(-l.offset.X, -l.offset.Y, g.Width()-l.offset.X, g.Height()-l.offset.Y)
	gs := []layout.Guide{g}

	return layout.Guide{
		Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y),
	}, gs
}

func (l *layouter) Notify(f func()) comm.Id {
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

func (l *layouter) Unnotify(id comm.Id) {
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

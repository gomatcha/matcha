// Package segmentview implements a native segmented control.
package segmentview

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/segmentview"
	"gomatcha.io/matcha/view"
)

type View struct {
	view.Embed
	Enabled       bool
	Momentary     bool
	Titles        []string
	Value         int
	OnValueChange func(value int)
	PaintStyle    *paint.Style
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{Embed: ctx.NewEmbed(key), Enabled: true}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(29)
		s.WidthEqual(l.MaxGuide().Width())
		s.TopEqual(l.MaxGuide().Top())
		s.LeftEqual(l.MaxGuide().Left())
	})

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:        painter,
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/segmentview",
		NativeViewState: &segmentview.View{
			Value:     int64(v.Value),
			Titles:    v.Titles,
			Enabled:   v.Enabled,
			Momentary: v.Momentary,
		},
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				event := &segmentview.Event{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Value = int(event.Value)
				if v.OnValueChange != nil {
					v.OnValueChange(v.Value)
				}
			},
		},
	}
}

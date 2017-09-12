package ios

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	pbios "gomatcha.io/matcha/pb/view/ios"
	"gomatcha.io/matcha/view"
)

type SegmentView struct {
	view.Embed
	Enabled       bool
	Momentary     bool
	Titles        []string
	Value         int
	OnValueChange func(value int)
	PaintStyle    *paint.Style
}

// NewSegmentView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewSegmentView() *SegmentView {
	return &SegmentView{
		Enabled: true,
	}
}

// Build implements view.View.
func (v *SegmentView) Build(ctx *view.Context) view.Model {
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
		NativeViewState: &pbios.SegmentView{
			Value:     int64(v.Value),
			Titles:    v.Titles,
			Enabled:   v.Enabled,
			Momentary: v.Momentary,
		},
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				event := &pbios.SegmentViewEvent{}
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

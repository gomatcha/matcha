// Package segmentview implements a native slider.
package slider

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/slider"
	"gomatcha.io/matcha/view"
)

type layouter struct {
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize.X, 31)}
	return g, nil
}

func (l *layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *layouter) Unnotify(id comm.Id) {
	// no-op
}

type View struct {
	view.Embed
	PaintStyle    *paint.Style
	Value         float64
	ValueNotifier comm.Float64Notifier
	valueNotifier comm.Float64Notifier
	MaxValue      float64
	MinValue      float64
	OnValueChange func(value float64)
	OnSubmit      func(value float64)
	Enabled       bool
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed:    ctx.NewEmbed(key),
		MaxValue: 1,
		MinValue: 0,
		Enabled:  true,
	}
}

func (v *View) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		if v.valueNotifier != nil {
			v.Unsubscribe(v.valueNotifier)
		}
	}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	val := v.Value
	if v.ValueNotifier != nil {
		val = v.ValueNotifier.Value()
	}

	if v.ValueNotifier != v.valueNotifier {
		if v.valueNotifier != nil {
			v.Unsubscribe(v.valueNotifier)
		}
		if v.ValueNotifier != nil {
			v.Subscribe(v.ValueNotifier)
		}
		v.valueNotifier = v.ValueNotifier
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:        painter,
		Layouter:       &layouter{},
		NativeViewName: "gomatcha.io/matcha/view/slider",
		NativeViewState: &slider.View{
			Value:    val,
			MaxValue: v.MaxValue,
			MinValue: v.MinValue,
			Enabled:  v.Enabled,
		},
		NativeFuncs: map[string]interface{}{
			"OnValueChange": func(data []byte) {
				event := &slider.Event{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				if v.OnValueChange != nil {
					v.OnValueChange(event.Value)
				}
			},
			"OnSubmit": func(data []byte) {
				event := &slider.Event{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				if v.OnSubmit != nil {
					v.OnSubmit(event.Value)
				}
			},
		},
	}
}

package view

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pbview "gomatcha.io/matcha/pb/view"
)

type sliderLayouter struct {
}

func (l *sliderLayouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize().X, 31)}
	return g, nil
}

func (l *sliderLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *sliderLayouter) Unnotify(id comm.Id) {
	// no-op
}

type Slider struct {
	Embed
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
func NewSlider() *Slider {
	return &Slider{
		MaxValue: 1,
		MinValue: 0,
		Enabled:  true,
	}
}

func (v *Slider) Lifecycle(from, to Stage) {
	if ExitsStage(from, to, StageMounted) {
		if v.valueNotifier != nil {
			v.Unsubscribe(v.valueNotifier)
		}
	}
}

// Build implements view.View.
func (v *Slider) Build(ctx Context) Model {
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
	return Model{
		Painter:        painter,
		Layouter:       &sliderLayouter{},
		NativeViewName: "gomatcha.io/matcha/view/slider",
		NativeViewState: &pbview.Slider{
			Value:    val,
			MaxValue: v.MaxValue,
			MinValue: v.MinValue,
			Enabled:  v.Enabled,
		},
		NativeFuncs: map[string]interface{}{
			"OnValueChange": func(data []byte) {
				event := &pbview.SliderEvent{}
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
				event := &pbview.SliderEvent{}
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

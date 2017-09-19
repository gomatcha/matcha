package view

import (
	"fmt"
	"runtime"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pbview "gomatcha.io/matcha/pb/view"
)

type Switch struct {
	Embed
	Enabled       bool
	Value         bool
	OnValueChange func(value bool)
	PaintStyle    *paint.Style
}

// NewSwitch returns either the previous View in ctx with matching key, or a new View if none exists.
func NewSwitch() *Switch {
	return &Switch{
		Enabled: true,
	}
}

// Build implements view.View.
func (v *Switch) Build(ctx Context) Model {
	var rect layout.Rect
	if runtime.GOOS == "android" {
		rect = layout.Rt(0, 0, 61, 40)
	} else {
		rect = layout.Rt(0, 0, 51, 31)
	}

	l := &absoluteLayouter{Guide: layout.Guide{Frame: rect}}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Painter:        painter,
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/switch",
		NativeViewState: &pbview.SwitchView{
			Value:   v.Value,
			Enabled: v.Enabled,
		},
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				event := &pbview.SwitchEvent{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Value = event.Value
				if v.OnValueChange != nil {
					v.OnValueChange(v.Value)
				}
			},
		},
	}
}

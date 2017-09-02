package view

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/switchview"
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
func (v *Switch) Build(ctx *Context) Model {
	l := &absoluteLayouter{
		Guide: layout.Guide{
			Frame: layout.Rt(0, 0, 51, 31),
		},
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Painter:        painter,
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/switch",
		NativeViewState: &switchview.View{
			Value:   v.Value,
			Enabled: v.Enabled,
		},
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				event := &switchview.Event{}
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

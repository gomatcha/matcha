package view

import (
	"fmt"
	"runtime"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	protoview "gomatcha.io/matcha/proto/view"
)

type Switch struct {
	Embed
	Enabled    bool
	Value      bool
	OnSubmit   func(value bool)
	PaintStyle *paint.Style
}

// NewSwitch returns a new view.
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
		NativeViewState: internal.MarshalProtobuf(&protoview.SwitchView{
			Value:   v.Value,
			Enabled: v.Enabled,
		}),
		NativeFuncs: map[string]interface{}{
			// "test": func(a, b string) {
			// 	fmt.Println("test", a, b)
			// },
			"OnChange": func(data []byte) {
				event := &protoview.SwitchEvent{}
				err := proto.Unmarshal(data, event)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.Value = event.Value
				if v.OnSubmit != nil {
					v.OnSubmit(v.Value)
				}
			},
		},
	}
}

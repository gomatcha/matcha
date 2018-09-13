package customview

import (
	"fmt"
	"runtime"

	"github.com/gogo/protobuf/proto"
	protoview "github.com/gomatcha/matcha/examples/customview/proto"
	"github.com/gomatcha/matcha/internal"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/view"
)

type CustomView struct {
	view.Embed
	Enabled  bool
	Value    bool
	OnSubmit func(value bool)
}

// NewCustomView returns an initialized CustomView instance.
func NewCustomView() *CustomView {
	return &CustomView{
		Enabled: true,
	}
}

// Build implements view.View.
func (v *CustomView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		if runtime.GOOS == "android" {
			s.Width(61)
			s.Height(40)
		} else {
			s.Width(51)
			s.Height(31)
		}
	})
	return view.Model{
		Layouter:       l,
		NativeViewName: "github.com/gomatcha/matcha/view/switch",
		NativeViewState: internal.MarshalProtobuf(&protoview.View{
			Value:   v.Value,
			Enabled: v.Enabled,
		}),
		NativeFuncs: map[string]interface{}{
			"OnChange": func(data []byte) {
				event := &protoview.Event{}
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

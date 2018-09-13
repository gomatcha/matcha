package bridge

import (
	"fmt"
	"runtime"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/bridge NewBridgeView", func() view.View {
		return NewBridgeView()
	})
}

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/bridge callWithGoValues", func(v int64) string {
		return fmt.Sprintf("Call with Go values:%v", v)
	})
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/bridge callWithForeignValues", func(v *bridge.Value) *bridge.Value {
		return bridge.String(fmt.Sprintf("Call with Foreign values:%v", v.ToInt64()))
	})
}

type BridgeView struct {
	view.Embed
}

func NewBridgeView() *BridgeView {
	return &BridgeView{}
}

func (v *BridgeView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	var str string
	var str2 string
	if runtime.GOOS == "android" {
		str = bridge.Bridge("github.com/gomatcha/matcha/example").Call("callWithGoValues", bridge.Interface(123)).ToInterface().(string)
		str2 = bridge.Bridge("github.com/gomatcha/matcha/example").Call("callWithForeignValues", bridge.Int64(456)).ToString()
	} else {
		str = bridge.Bridge("github.com/gomatcha/matcha/example").Call("callWithGoValues:", bridge.Interface(123)).ToInterface().(string)
		str2 = bridge.Bridge("github.com/gomatcha/matcha/example").Call("callWithForeignValues:", bridge.Int64(456)).ToString()
	}

	chl1 := view.NewTextView()
	chl1.String = str
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.TopEqual(l.Top().Add(50))
		s.LeftEqual(l.Left())
	})

	chl2 := view.NewTextView()
	chl2.String = str2
	_ = l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

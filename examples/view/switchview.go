package view

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewSwitchView", func() view.View {
		return NewSwitchView()
	})
}

type SwitchView struct {
	view.Embed
	value bool
}

func NewSwitchView() *SwitchView {
	return &SwitchView{
		value: true,
	}
}

func (v *SwitchView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewSwitch()
	chl1.Value = v.value
	chl1.PaintStyle = &paint.Style{BackgroundColor: colornames.Green}
	chl1.OnSubmit = func(value bool) {
		v.value = value
		v.Signal()
		fmt.Println("onValueChange", value)
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
	})

	chl2 := view.NewSwitch()
	chl2.Value = v.value
	chl2.Enabled = false
	chl2.OnSubmit = func(value bool) {
		fmt.Println("onValueChange2", value)
	}
	l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(50))
		s.LeftEqual(g1.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

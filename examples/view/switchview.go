package view

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/basicview"
	"gomatcha.io/matcha/view/switchview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewSwitchView", func() *view.Root {
		return view.NewRoot(NewSwitchView())
	})
}

type SwitchView struct {
	view.Embed
	value *comm.Float64Value
}

func NewSwitchView() *SwitchView {
	return &SwitchView{
		value: &comm.Float64Value{},
	}
}

func (v *SwitchView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl0 := basicview.New()
	chl0.Painter = &paint.Style{BackgroundColor: colornames.Red}
	l.Add(chl0, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(100)
		s.Height(100)
	})

	chl1 := switchview.New()
	chl1.Value = true
	chl1.PaintStyle = &paint.Style{BackgroundColor: colornames.Green}
	chl1.OnValueChange = func(value bool) {
		fmt.Println("onValueChange", value)
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
	})

	chl2 := switchview.New()
	chl2.Value = false
	chl2.OnValueChange = func(value bool) {
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

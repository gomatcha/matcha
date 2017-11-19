package view

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
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

	label := view.NewTextView()
	label.String = "Enabled Switch:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	s := view.NewSwitch()
	s.Value = v.value
	s.OnSubmit = func(value bool) {
		v.value = value
		v.Signal()
		fmt.Println("OnSubmit", value)
	}
	g = l.Add(s, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Disabled Switch:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	s = view.NewSwitch()
	s.Value = v.value
	s.Enabled = false
	s.OnSubmit = func(value bool) {
		fmt.Println("OnSubmit", value)
	}
	l.Add(s, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

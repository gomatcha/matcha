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
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view Slider", func() view.View {
		return NewSlider()
	})
}

type Slider struct {
	view.Embed
	value float64
}

func NewSlider() *Slider {
	return &Slider{
		value: 0.5,
	}
}

func (v *Slider) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := view.NewSlider()
	chl1.MinValue = -5
	chl1.MaxValue = 5
	chl1.Value = v.value
	chl1.PaintStyle = &paint.Style{BackgroundColor: colornames.Green}
	chl1.OnValueChange = func(value float64) {
		v.value = value
		v.Signal()
		fmt.Println("onValueChange", value)
	}
	chl1.OnSubmit = func(value float64) {
		fmt.Println("onSubmit", value)
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := view.NewSlider()
	chl2.Value = v.value
	chl2.MinValue = -10
	chl2.MaxValue = 10
	chl2.Enabled = false
	chl2.OnValueChange = func(value float64) {
		fmt.Println("onValueChange2", value)
	}
	l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(50))
		s.LeftEqual(g1.Left())
		s.WidthEqual(g1.Width())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

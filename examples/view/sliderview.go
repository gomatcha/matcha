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
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewSliderView", func() view.View {
		return NewSliderView()
	})
}

type SliderView struct {
	view.Embed
	value float64
}

func NewSliderView() *SliderView {
	return &SliderView{
		value: 0.5,
	}
}

func (v *SliderView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Enabled Slider:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	slider := view.NewSlider()
	slider.MinValue = -5
	slider.MaxValue = 5
	slider.Value = v.value
	slider.OnChange = func(value float64) {
		v.value = value
		v.Signal()
		fmt.Println("onValueChange", value)
	}
	slider.OnSubmit = func(value float64) {
		fmt.Println("onSubmit", value)
	}
	g = l.Add(slider, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	label = view.NewTextView()
	label.String = "Disabled Slider:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	slider = view.NewSlider()
	slider.Value = v.value
	slider.MinValue = 0
	slider.MaxValue = 5
	slider.Enabled = false
	slider.OnChange = func(value float64) {
		fmt.Println("onValueChange2", value)
	}
	g = l.Add(slider, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

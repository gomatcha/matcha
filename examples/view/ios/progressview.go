package ios

import (
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"github.com/gomatcha/matcha/view/ios"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view NewProgressView", func() view.View {
		return NewProgressView()
	})
}

type ProgressView struct {
	view.Embed
	value *comm.Float64Value
}

func NewProgressView() *ProgressView {
	return &ProgressView{
		value: &comm.Float64Value{},
	}
}

func (v *ProgressView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	progressv := ios.NewProgressView()
	progressv.ProgressNotifier = v.value
	progressv.ProgressColor = colornames.Red
	l.Add(progressv, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	sliderv := view.NewSlider()
	sliderv.MaxValue = 1
	sliderv.MinValue = 0
	sliderv.OnChange = func(value float64) {
		v.value.SetValue(value)
	}
	l.Add(sliderv, func(s *constraint.Solver) {
		s.Top(200)
		s.Left(100)
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

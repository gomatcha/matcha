package ios

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewProgressView", func() view.View {
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

	label := view.NewTextView()
	label.String = "Progress view:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	progress := ios.NewProgressView()
	progress.ProgressNotifier = v.value
	progress.ProgressColor = colornames.Red
	g = l.Add(progress, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom().Add(10))
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	label = view.NewTextView()
	label.String = "Progress value:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom().Add(10))
		s.LeftEqual(g.Left())
	})

	slider := view.NewSlider()
	slider.MaxValue = 1
	slider.MinValue = 0
	slider.OnChange = func(value float64) {
		v.value.SetValue(value)
	}
	l.Add(slider, func(s *constraint.Solver) {
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

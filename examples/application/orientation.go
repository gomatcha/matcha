package application

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/application NewOrientationView", func() view.View {
		return NewOrientationView()
	})
}

type OrientationView struct {
	view.Embed
	notifier comm.IntNotifier
}

func NewOrientationView() *OrientationView {
	return &OrientationView{
		notifier: application.OrientationNotifier,
	}
}

func (v *OrientationView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.notifier)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.notifier)
	}
}

func (v *OrientationView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Orientation:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	var str string
	switch layout.Edge(v.notifier.Value()) {
	case layout.EdgeTop:
		str = "Top"
	case layout.EdgeBottom:
		str = "Bottom"
	case layout.EdgeRight:
		str = "Right"
	case layout.EdgeLeft:
		str = "Left"
	}

	label = view.NewTextView()
	label.String = str
	label.Style.SetFont(text.DefaultFont(15))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Get Orientation"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button := view.NewButton()
	button.String = "Button"
	button.OnPress = func() {
		var str string
		switch application.Orientation() {
		case layout.EdgeTop:
			str = "Top"
		case layout.EdgeBottom:
			str = "Bottom"
		case layout.EdgeRight:
			str = "Right"
		case layout.EdgeLeft:
			str = "Left"
		}
		view.Alert("Orientation", str)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

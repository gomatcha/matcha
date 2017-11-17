package ios

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/ios NewActivityIndicatorView", func() view.View {
		return NewActivityIndicatorView()
	})
}

type ActivityIndicatorView struct {
	view.Embed
	hidden bool
}

func NewActivityIndicatorView() *ActivityIndicatorView {
	return &ActivityIndicatorView{}
}

func (v *ActivityIndicatorView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Toggle Activity Indicator:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	chl1 := view.NewButton()
	if v.hidden {
		chl1.String = "Show"
	} else {
		chl1.String = "Hide"
	}
	chl1.OnPress = func() {
		v.hidden = !v.hidden
		v.Signal()
	}
	l.Add(chl1, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	options := []view.Option{}
	if !v.hidden {
		options = append(options, &ios.ActivityIndicator{})
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options:  options,
	}
}

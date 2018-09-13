package ios

import (
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"github.com/gomatcha/matcha/view/ios"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/ios NewActivityIndicatorView", func() view.View {
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

	chl1 := view.NewButton()
	chl1.String = "Toggle Hidden"
	chl1.OnPress = func() {
		v.hidden = !v.hidden
		v.Signal()
	}
	l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	options := []view.Option{}
	if !v.hidden {
		options = append(options, &ios.ActivityIndicator{})
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Lightgray},
		Options:  options,
	}
}

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
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewStatusBarView", func() view.View {
		return NewStatusBarView()
	})
}

type StatusBarView struct {
	view.Embed
	style  ios.StatusBarStyle
	hidden bool
}

func NewStatusBarView() *StatusBarView {
	return &StatusBarView{
		style: ios.StatusBarStyleDark,
	}
}

func (v *StatusBarView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Toggle Style:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	button := view.NewButton()
	if v.style == ios.StatusBarStyleLight {
		button.String = "Dark"
	} else {
		button.String = "Light"
	}
	button.OnPress = func() {
		if v.style == ios.StatusBarStyleLight {
			v.style = ios.StatusBarStyleDark
		} else {
			v.style = ios.StatusBarStyleLight
		}
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Toggle Hidden:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	if v.hidden {
		button.String = "Show"
	} else {
		button.String = "Hide"
	}
	button.OnPress = func() {
		v.hidden = !v.hidden
		v.Signal()
	}
	l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			&ios.StatusBar{
				Style:  v.style,
				Hidden: v.hidden,
			},
		},
	}
}

package android

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"

	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/android NewPagerView", func() view.View {
		return NewPagerView(NewPagerApp())
	})
}

type PagerApp struct {
	pages          *android.Pages
	relay          *comm.Relay
	barColor       color.Color
	indicatorColor color.Color
}

func NewPagerApp() *PagerApp {
	app := &PagerApp{}
	app.pages = &android.Pages{}
	app.pages.Notify(func() {
		fmt.Println("CurrentPage", app.pages.SelectedIndex())
	})

	v1 := NewPagerChild(app)
	v1.Color = colornames.Yellow
	v1.Title = "Search"

	v2 := NewPagerChild(app)
	v2.Color = colornames.Red
	v2.Title = "Settings"

	v3 := NewPagerChild(app)
	v3.Color = colornames.White
	v3.Title = "Map"

	v4 := NewPagerChild(app)
	v4.Color = colornames.Green
	v4.Title = "Camera"

	app.pages.SetViews(v1, v2, v3, v4)
	app.pages.SetSelectedIndex(2)
	app.relay = &comm.Relay{}
	return app
}

type PagerView struct {
	view.Embed
	app *PagerApp
}

func NewPagerView(app *PagerApp) view.View {
	return &PagerView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *PagerView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.app.relay)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.app.relay)
	}
}

func (v *PagerView) Build(ctx view.Context) view.Model {
	pagerview := android.NewPagerView()
	pagerview.Pages = v.app.pages
	// pagerview.BarColor = v.app.barColor
	// pagerview.IndicatorColor = v.app.indicatorColor
	return view.Model{
		Children: []view.View{pagerview},
	}
}

type PagerChild struct {
	view.Embed
	app            *PagerApp
	Color          color.Color
	Title          string
	indicatorColor color.Color
}

func NewPagerChild(app *PagerApp) *PagerChild {
	return &PagerChild{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *PagerChild) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Current page:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	button := view.NewButton()
	button.String = "0"
	button.OnPress = func() {
		v.app.pages.SetSelectedIndex(0)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "1"
	button.OnPress = func() {
		v.app.pages.SetSelectedIndex(1)
	}
	buttonG := l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Top())
		s.LeftEqual(g.Right())
	})

	button = view.NewButton()
	button.String = "2"
	button.OnPress = func() {
		v.app.pages.SetSelectedIndex(2)
	}
	buttonG = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(buttonG.Top())
		s.LeftEqual(buttonG.Right())
	})

	button = view.NewButton()
	button.String = "3"
	button.OnPress = func() {
		v.app.pages.SetSelectedIndex(3)
	}
	buttonG = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(buttonG.Top())
		s.LeftEqual(buttonG.Right())
	})

	label = view.NewTextView()
	label.String = "Bar Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.barColor == nil {
			v.app.barColor = colornames.Yellow
		} else {
			v.app.barColor = nil
		}
		v.app.relay.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "All Indicator Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.indicatorColor == nil {
			v.app.indicatorColor = colornames.Teal
		} else {
			v.app.indicatorColor = nil
		}
		v.app.relay.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Indicator Color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.indicatorColor == nil {
			v.indicatorColor = colornames.Red
		} else {
			v.indicatorColor = nil
		}
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: v.Color},
		Options: []view.Option{
			&android.PagerButton{
				Title: v.Title,
				// IndicatorColor: v.indicatorColor,
			},
		},
	}
}

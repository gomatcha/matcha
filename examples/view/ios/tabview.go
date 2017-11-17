package ios

import (
	"image"
	"image/color"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewTabView", func() view.View {
		return NewTabView(NewTabApp())
	})
}

type TabApp struct {
	tabs *ios.Tabs

	relay           *comm.Relay
	barColor        color.Color
	selectedColor   color.Color
	unselectedColor color.Color
}

func NewTabApp() *TabApp {
	app := &TabApp{
		tabs:  &ios.Tabs{},
		relay: &comm.Relay{},
	}

	view1 := NewTabChild(app)
	view1.Color = colornames.White
	view1.Title = "Search"
	view1.Icon = application.MustLoadImage("tab_search")
	view1.SelectedIcon = application.MustLoadImage("tab_search_filled")

	view2 := NewTabChild(app)
	view2.Color = colornames.Red
	view2.Title = "Settings"
	view2.Icon = application.MustLoadImage("tab_settings")
	view2.SelectedIcon = application.MustLoadImage("tab_settings_filled")

	view3 := NewTabChild(app)
	view3.Color = colornames.Yellow
	view3.Title = "Map"
	view3.Icon = application.MustLoadImage("tab_map")
	view3.SelectedIcon = application.MustLoadImage("tab_map_filled")

	view4 := NewTabChild(app)
	view4.Color = colornames.Green
	view4.Title = "Camera"
	view4.Icon = application.MustLoadImage("tab_camera")
	view4.SelectedIcon = application.MustLoadImage("tab_camera_filled")

	app.tabs.SetViews(view1, view2, view3, view4)
	app.tabs.SetSelectedIndex(1)
	return app
}

type TabView struct {
	view.Embed
	app *TabApp
}

func NewTabView(app *TabApp) view.View {
	return &TabView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *TabView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.app.relay)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.app.relay)
	}
}

func (v *TabView) Build(ctx view.Context) view.Model {
	tabview := ios.NewTabView()
	tabview.BarColor = v.app.barColor
	tabview.SelectedColor = v.app.selectedColor
	tabview.UnselectedColor = v.app.unselectedColor
	tabview.Tabs = v.app.tabs
	return view.Model{
		Children: []view.View{tabview},
	}
}

type TabChild struct {
	view.Embed
	app          *TabApp
	Color        color.Color
	Title        string
	Icon         image.Image
	SelectedIcon image.Image
	titleToggle  bool
	iconToggle   bool
	badgeToggle  bool
}

func NewTabChild(app *TabApp) *TabChild {
	return &TabChild{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *TabChild) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Current tab:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	button := view.NewButton()
	button.String = "0"
	button.OnPress = func() {
		v.app.tabs.SetSelectedIndex(0)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "1"
	button.OnPress = func() {
		v.app.tabs.SetSelectedIndex(1)
	}
	buttonG := l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Top())
		s.LeftEqual(g.Right())
	})

	button = view.NewButton()
	button.String = "2"
	button.OnPress = func() {
		v.app.tabs.SetSelectedIndex(2)
	}
	buttonG = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(buttonG.Top())
		s.LeftEqual(buttonG.Right())
	})

	button = view.NewButton()
	button.String = "3"
	button.OnPress = func() {
		v.app.tabs.SetSelectedIndex(3)
	}
	buttonG = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(buttonG.Top())
		s.LeftEqual(buttonG.Right())
	})

	label = view.NewTextView()
	label.String = "Tab title:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		v.titleToggle = !v.titleToggle
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Tab badge:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		v.badgeToggle = !v.badgeToggle
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Tab icon:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		v.iconToggle = !v.iconToggle
		v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Bar color:"
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
	label.String = "Selected color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.selectedColor == nil {
			v.app.selectedColor = colornames.Red
		} else {
			v.app.selectedColor = nil
		}
		v.app.relay.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Unselected color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Toggle"
	button.OnPress = func() {
		if v.app.unselectedColor == nil {
			v.app.unselectedColor = colornames.Green
		} else {
			v.app.unselectedColor = nil
		}
		v.app.relay.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	tab := &ios.TabButton{}
	if !v.titleToggle {
		tab.Title = v.Title
	} else {
		tab.Title = ""
	}
	if !v.iconToggle {
		tab.Icon = v.Icon
		tab.SelectedIcon = v.SelectedIcon
	}
	if v.badgeToggle {
		tab.Badge = "Badge"
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: v.Color},
		Options:  []view.Option{tab},
	}
}

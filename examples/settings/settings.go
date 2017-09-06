// Package settings provides a skeleton implemention the iOS settings UI.
package settings

import (
	"fmt"
	"image/color"
	"strings"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

var (
	cellColor            = color.Gray{255}
	cellColorHighlighted = color.Gray{217}
	chevronColor         = color.RGBA{199, 199, 204, 255}
	separatorColor       = color.RGBA{203, 202, 207, 255}
	backgroundColor      = color.RGBA{239, 239, 244, 255}
	subtitleColor        = color.Gray{142}
	titleColor           = color.Gray{0}
	spacerTitleColor     = color.Gray{102}
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/settings New", func() *view.Root {
		app := NewApp()

		v := ios.NewStackView()
		v.Stack = app.Stack
		v.Stack.SetViews(NewRootView(app))
		return view.NewRoot(v)
	})
}

type RootView struct {
	view.Embed
	app *App
}

func NewRootView(app *App) *RootView {
	return &RootView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *RootView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.app)
		v.Subscribe(v.app.Wifi)
		v.Subscribe(v.app.Bluetooth)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.app)
		v.Unsubscribe(v.app.Wifi)
		v.Unsubscribe(v.app.Bluetooth)
	}
}

func (v *RootView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}
	{
		group := []view.View{}

		spacer := NewSpacer()
		l.Add(spacer, nil)

		switchView := view.NewSwitch()
		switchView.Value = v.app.AirplaneMode()
		switchView.OnValueChange = func(value bool) {
			fmt.Println("blah")
			v.app.SetAirplaneMode(value)
		}
		cell1 := NewBasicCell()
		cell1.Title = "Airplane Mode"
		cell1.Icon = app.MustLoadImage("airplane")
		cell1.AccessoryView = switchView
		cell1.HasIcon = true
		group = append(group, cell1)

		cell2 := NewBasicCell()
		cell2.Title = "Wi-Fi"
		if v.app.Wifi.Enabled() {
			cell2.Subtitle = v.app.Wifi.CurrentSSID()
		} else {
			cell2.Subtitle = ""
		}
		fmt.Println("wifi", cell2.Subtitle)
		cell2.HasIcon = true
		cell2.Icon = app.MustLoadImage("Wifi")
		cell2.Chevron = true
		cell2.OnTap = func() {
			v.app.Stack.Push(NewWifiView(v.app))
		}
		group = append(group, cell2)

		cell3 := NewBasicCell()
		cell3.HasIcon = true
		cell3.Icon = app.MustLoadImage("Bluetooth")
		cell3.Title = "Bluetooth"
		if v.app.Bluetooth.Enabled() {
			cell3.Subtitle = "On"
		} else {
			cell3.Subtitle = ""
		}
		cell3.Chevron = true
		cell3.OnTap = func() {
			v.app.Stack.Push(NewBluetoothView(v.app))
		}
		group = append(group, cell3)

		cell4 := NewBasicCell()
		cell4.HasIcon = true
		cell4.Icon = app.MustLoadImage("Cellular")
		cell4.Title = "Cellular"
		cell4.Chevron = true
		cell4.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell4)

		cell5 := NewBasicCell()
		cell5.HasIcon = true
		cell5.Icon = app.MustLoadImage("Hotspot")
		cell5.Title = "Personal Hotspot"
		cell5.Subtitle = "Off"
		cell5.Chevron = true
		cell5.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell5)

		cell6 := NewBasicCell()
		cell6.HasIcon = true
		cell6.Icon = app.MustLoadImage("Carrier")
		cell6.Title = "Carrier"
		cell6.Subtitle = "T-Mobile"
		cell6.Chevron = true
		cell6.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell6)

		for _, i := range AddSeparators(group) {
			l.Add(i, nil)
		}
	}
	{
		group := []view.View{}

		spacer := NewSpacer()
		l.Add(spacer, nil)

		cell1 := NewBasicCell()
		cell1.HasIcon = true
		cell1.Icon = app.MustLoadImage("Notifications")
		cell1.Title = "Notifications"
		cell1.Chevron = true
		cell1.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell1)

		cell2 := NewBasicCell()
		cell2.HasIcon = true
		cell2.Icon = app.MustLoadImage("ControlCenter")
		cell2.Title = "Control Center"
		cell2.Chevron = true
		cell2.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell2)

		cell3 := NewBasicCell()
		cell3.HasIcon = true
		cell3.Icon = app.MustLoadImage("DoNotDisturb")
		cell3.Title = "Do Not Disturb"
		cell3.Chevron = true
		cell3.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell3)

		for _, i := range AddSeparators(group) {
			l.Add(i, nil)
		}
	}

	scrollView := view.NewScrollView()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

func (v *RootView) StackBar(ctx *view.Context) *ios.Bar {
	return &ios.Bar{Title: "Settings Example"}
}

func AddSeparators(vs []view.View) []view.View {
	newViews := []view.View{}

	top := NewSeparator()
	newViews = append(newViews, top)

	for idx, i := range vs {
		newViews = append(newViews, i)

		if idx != len(vs)-1 { // Don't add short separator after last view
			sep := NewSeparator()
			sep.LeftPadding = 60
			newViews = append(newViews, sep)
		}
	}

	bot := NewSeparator()
	newViews = append(newViews, bot)
	return newViews
}

type Separator struct {
	view.Embed
	LeftPadding float64
}

func NewSeparator() *Separator {
	return &Separator{}
}

func (v *Separator) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(0.5)
		s.WidthEqual(l.MaxGuide().Width())
	})

	chl := view.NewBasicView()
	chl.Painter = &paint.Style{BackgroundColor: separatorColor}
	l.Add(chl, func(s *constraint.Solver) {
		s.HeightEqual(l.Height())
		s.LeftEqual(l.Left().Add(v.LeftPadding))
		s.RightEqual(l.Right())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: cellColor},
	}
}

type Spacer struct {
	view.Embed
	Height float64
}

func NewSpacer() *Spacer {
	return &Spacer{
		Height: 35,
	}
}

func (v *Spacer) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(v.Height)
		s.WidthEqual(l.MaxGuide().Width())
	})

	return view.Model{
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

type SpacerHeader struct {
	view.Embed
	Height float64
	Title  string
}

func NewSpacerHeader() *SpacerHeader {
	return &SpacerHeader{
		Height: 50,
	}
}

func (v *SpacerHeader) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(v.Height)
		s.WidthEqual(l.MaxGuide().Width())
	})

	titleView := view.NewTextView()
	titleView.String = strings.ToTitle(v.Title)
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 13))
	titleView.Style.SetTextColor(spacerTitleColor)
	// titleView.Painter = &paint.Style{BackgroundColor: colornames.Red}

	titleGuide := l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(15))
		s.RightEqual(l.Right().Add(-15))
		s.BottomEqual(l.Bottom().Add(-10))
		s.TopGreater(l.Top())
	})
	_ = titleGuide

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

type SpacerDescription struct {
	view.Embed
	Description string
}

func NewSpacerDescription() *SpacerDescription {
	return &SpacerDescription{}
}

func (v *SpacerDescription) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	titleView := view.NewTextView()
	titleView.String = v.Description
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 13))
	titleView.Style.SetTextColor(spacerTitleColor)

	titleGuide := l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(15))
		s.RightEqual(l.Right().Add(-15))
		s.TopGreater(l.Top().Add(15))
	})

	l.Solve(func(s *constraint.Solver) {
		s.HeightEqual(titleGuide.Height().Add(30))
		s.WidthEqual(l.MaxGuide().Width())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

type BasicCell struct {
	view.Embed
	HasIcon       bool
	Icon          *app.ImageResource
	Title         string
	Subtitle      string
	AccessoryView view.View
	Chevron       bool
	OnTap         func()
	highlighted   bool
}

func NewBasicCell() *BasicCell {
	return &BasicCell{}
}

func (v *BasicCell) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(44)
		s.WidthEqual(l.MaxGuide().Width())
	})

	leftAnchor := l.Left()
	if v.HasIcon {
		iconView := view.NewImageView()
		iconView.Image = v.Icon
		iconView.ResizeMode = view.ImageResizeModeFill
		pIconView := view.WithPainter(iconView, &paint.Style{BackgroundColor: colornames.Lightgray, CornerRadius: 5})

		iconGuide := l.Add(pIconView, func(s *constraint.Solver) {
			s.Width(30)
			s.Height(30)
			s.LeftEqual(l.Left().Add(15))
			s.CenterYEqual(l.CenterY())
		})
		leftAnchor = iconGuide.Right()
	}

	rightAnchor := l.Right()
	if v.Chevron {
		chevronView := view.NewImageView()
		chevronView.Image = app.MustLoadImage("table_arrow")
		chevronView.ResizeMode = view.ImageResizeModeCenter
		chevronView.ImageTint = chevronColor

		chevronGuide := l.Add(chevronView, func(s *constraint.Solver) {
			s.RightEqual(rightAnchor.Add(-15))
			s.LeftGreater(leftAnchor)
			s.CenterYEqual(l.CenterY())
			s.TopGreater(l.Top())
			s.BottomLess(l.Bottom())
		})
		rightAnchor = chevronGuide.Left()
	}

	if v.AccessoryView != nil {
		accessoryGuide := l.Add(v.AccessoryView, func(s *constraint.Solver) {
			s.RightEqual(rightAnchor.Add(-10))
			s.LeftGreater(leftAnchor)
			s.CenterYEqual(l.CenterY())
		})
		rightAnchor = accessoryGuide.Left()
	}

	if len(v.Subtitle) > 0 {
		fmt.Println("subtitle", v.Subtitle)
		subtitleView := view.NewTextView()
		subtitleView.String = v.Subtitle
		subtitleView.Style.SetFont(text.FontWithName("HelveticaNeue", 14))
		subtitleView.Style.SetTextColor(subtitleColor)

		subtitleGuide := l.Add(subtitleView, func(s *constraint.Solver) {
			s.RightEqual(rightAnchor.Add(-10))
			s.LeftGreater(leftAnchor)
			s.CenterYEqual(l.CenterY())
		})
		rightAnchor = subtitleGuide.Left()
	}

	titleView := view.NewTextView()
	titleView.String = v.Title
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 14))
	titleView.Style.SetTextColor(titleColor)

	titleGuide := l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(leftAnchor.Add(15))
		s.RightLess(rightAnchor.Add(-10))
		s.CenterYEqual(l.CenterY())
	})
	_ = titleGuide

	var options []view.Option
	if v.OnTap != nil {
		tap := &touch.ButtonGesture{
			OnTouch: func(e *touch.ButtonEvent) {
				switch e.Kind {
				case touch.EventKindPossible:
					v.highlighted = e.Inside
				case touch.EventKindFailed:
					v.highlighted = false
				case touch.EventKindRecognized:
					v.highlighted = false
					v.OnTap()
				}
				v.Signal()
			},
		}
		options = append(options, touch.GestureList{tap})
	}

	var color color.Color
	if v.highlighted {
		color = cellColorHighlighted
	} else {
		color = cellColor
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: color},
		Options:  options,
	}
}

// Package settings provides a skeleton implemention the iOS settings UI.
package settings

import (
	"image/color"
	"runtime"
	"strings"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/examples/internal"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
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
	BackgroundColor      = backgroundColor
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/settings New", func() view.View {
		return NewRootView()
	})
}

func NewRootView() view.View {
	if runtime.GOOS == "android" {
		v := android.NewStackView()
		app := NewApp()
		app.Stack = v.Stack
		app.Stack.SetViews(NewAppView(app))
		return v
	} else {
		v := ios.NewStackView()
		app := NewApp()
		app.Stack = v.Stack
		app.Stack.SetViews(NewAppView(app))
		return v
	}
}

type AppView struct {
	view.Embed
	app *App
}

func NewAppView(app *App) *AppView {
	return &AppView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *AppView) Lifecycle(from, to view.Stage) {
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

func (v *AppView) Build(ctx view.Context) view.Model {
	l := &table.Layouter{}
	{
		group := []view.View{}

		spacer := NewSpacer()
		l.Add(spacer, nil)

		switchView := view.NewSwitch()
		switchView.Value = v.app.AirplaneMode()
		switchView.OnSubmit = func(value bool) {
			v.app.SetAirplaneMode(value)
		}
		cell1 := NewBasicCell()
		cell1.Title = "Airplane Mode"
		cell1.Icon = application.MustLoadImage("settings_airplane")
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
		cell2.HasIcon = true
		cell2.Icon = application.MustLoadImage("settings_wifi")
		cell2.Chevron = true
		cell2.OnTap = func() {
			v.app.Stack.Push(NewWifiView(v.app))
		}
		group = append(group, cell2)

		cell3 := NewBasicCell()
		cell3.HasIcon = true
		cell3.Icon = application.MustLoadImage("settings_bluetooth")
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
		cell4.Icon = application.MustLoadImage("settings_cellular")
		cell4.Title = "Cellular"
		cell4.Chevron = true
		cell4.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell4)

		cell5 := NewBasicCell()
		cell5.HasIcon = true
		cell5.Icon = application.MustLoadImage("settings_hotspot")
		cell5.Title = "Personal Hotspot"
		cell5.Subtitle = "Off"
		cell5.Chevron = true
		cell5.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell5)

		cell6 := NewBasicCell()
		cell6.HasIcon = true
		cell6.Icon = application.MustLoadImage("settings_carrier")
		cell6.Title = "Carrier"
		cell6.Subtitle = "T-Mobile"
		cell6.Chevron = true
		cell6.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell6)

		for _, i := range AddSeparators(group, 60) {
			l.Add(i, nil)
		}
	}
	{
		group := []view.View{}

		spacer := NewSpacer()
		l.Add(spacer, nil)

		cell1 := NewBasicCell()
		cell1.HasIcon = true
		cell1.Icon = application.MustLoadImage("settings_notifications")
		cell1.Title = "Notifications"
		cell1.Chevron = true
		cell1.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell1)

		cell2 := NewBasicCell()
		cell2.HasIcon = true
		cell2.Icon = application.MustLoadImage("settings_control_center")
		cell2.Title = "Control Center"
		cell2.Chevron = true
		cell2.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell2)

		cell3 := NewBasicCell()
		cell3.HasIcon = true
		cell3.Icon = application.MustLoadImage("settings_do_not_disturb")
		cell3.Title = "Do Not Disturb"
		cell3.Chevron = true
		cell3.OnTap = func() {
			v.app.Stack.Push(NewCellularView(v.app))
		}
		group = append(group, cell3)

		for _, i := range AddSeparators(group, 60) {
			l.Add(i, nil)
		}
	}

	scrollView := view.NewScrollView()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	iosStackBar := &ios.StackBar{Title: "Settings"}
	androidStackBar := &android.StackBar{Title: "Settings"}
	if item := internal.IosBackItem(); item != nil { // Add button to example list
		iosStackBar.LeftItems = []*ios.StackBarItem{item}
	}
	if item := internal.AndroidBackItem(); item != nil {
		androidStackBar.Items = []*android.StackBarItem{item}
	}

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
		Options: []view.Option{
			iosStackBar,
			androidStackBar,
		},
	}
}

func AddSeparators(vs []view.View, leftPadding float64) []view.View {
	newViews := []view.View{}

	top := NewSeparator()
	newViews = append(newViews, top)

	for idx, i := range vs {
		newViews = append(newViews, i)

		if idx != len(vs)-1 { // Don't add short separator after last view
			sep := NewSeparator()
			sep.LeftPadding = leftPadding
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

func (v *Separator) Build(ctx view.Context) view.Model {
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

func (v *Spacer) Build(ctx view.Context) view.Model {
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

func (v *SpacerHeader) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(v.Height)
		s.WidthEqual(l.MaxGuide().Width())
	})

	titleView := view.NewTextView()
	titleView.String = strings.ToTitle(v.Title)
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 13))
	titleView.Style.SetTextColor(spacerTitleColor)

	titleGuide := l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(15))
		s.RightEqual(l.Right().Add(-15))
		s.BottomEqual(l.Bottom().Add(-10))
		// s.TopGreater(l.Top()) // TODO(KD): Why does this affect the layout?
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

func (v *SpacerDescription) Build(ctx view.Context) view.Model {
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

type LargeCell struct {
	view.Embed
	Title       string
	OnTap       func()
	highlighted bool
}

func NewLargeCell() *LargeCell {
	return &LargeCell{}
}

func (v *LargeCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(60)
		s.WidthEqual(l.MaxGuide().Width())
	})

	chevronView := view.NewImageView()
	chevronView.Image = application.MustLoadImage("table_arrow")
	chevronView.ResizeMode = view.ImageResizeModeCenter
	chevronView.ImageTint = chevronColor

	chevronGuide := l.Add(chevronView, func(s *constraint.Solver) {
		s.Right(-15)
		s.CenterY(0)
	})

	titleView := view.NewTextView()
	titleView.String = v.Title
	titleView.Style.SetFont(text.FontWithName("HelveticaNeue", 24))
	titleView.Style.SetTextColor(colornames.Black)

	titleGuide := l.Add(titleView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(15))
		s.RightLess(chevronGuide.Left().Add(-15))
		s.CenterYEqual(l.CenterY())
	})
	_ = titleGuide

	var options []view.Option
	if v.OnTap != nil {
		tap := &pointer.ButtonGesture{
			OnEvent: func(e *pointer.ButtonEvent) {
				switch e.Kind {
				case pointer.EventKindPossible:
					v.highlighted = e.Inside
				case pointer.EventKindFailed:
					v.highlighted = false
				case pointer.EventKindRecognized:
					v.highlighted = false
					v.OnTap()
				}
				v.Signal()
			},
		}
		options = append(options, pointer.GestureList{tap})
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

type BasicCell struct {
	view.Embed
	HasIcon       bool
	Icon          *application.ImageResource
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

func (v *BasicCell) Build(ctx view.Context) view.Model {
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
		chevronView.Image = application.MustLoadImage("table_arrow")
		chevronView.ResizeMode = view.ImageResizeModeCenter
		chevronView.ImageTint = chevronColor

		chevronGuide := l.Add(chevronView, func(s *constraint.Solver) {
			s.Right(-15)
			s.CenterY(0)
		})
		rightAnchor = chevronGuide.Left()
	}

	if v.AccessoryView != nil {
		accessoryGuide := l.Add(v.AccessoryView, func(s *constraint.Solver) {
			s.RightEqual(rightAnchor.Add(-10))
			s.LeftGreater(leftAnchor)
			s.CenterY(0)
		})
		rightAnchor = accessoryGuide.Left()
	}

	if len(v.Subtitle) > 0 {
		subtitleView := view.NewTextView()
		subtitleView.String = v.Subtitle
		subtitleView.Style.SetFont(text.FontWithName("HelveticaNeue", 14))
		subtitleView.Style.SetTextColor(subtitleColor)

		subtitleGuide := l.Add(subtitleView, func(s *constraint.Solver) {
			s.RightEqual(rightAnchor.Add(-10))
			s.LeftGreater(leftAnchor)
			s.CenterY(0)
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
		tap := &pointer.ButtonGesture{
			OnEvent: func(e *pointer.ButtonEvent) {
				switch e.Kind {
				case pointer.EventKindPossible:
					v.highlighted = e.Inside
				case pointer.EventKindFailed:
					v.highlighted = false
				case pointer.EventKindRecognized:
					v.highlighted = false
					v.OnTap()
				}
				v.Signal()
			},
		}
		options = append(options, pointer.GestureList{tap})
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

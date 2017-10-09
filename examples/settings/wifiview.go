package settings

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

type WifiView struct {
	view.Embed
	app *App
}

func NewWifiView(app *App) *WifiView {
	return &WifiView{
		Embed: view.NewEmbed(app),
		app:   app,
	}
}

func (v *WifiView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.app.Wifi)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.app.Wifi)
	}
}

func (v *WifiView) Build(ctx view.Context) view.Model {
	l := &table.Layouter{}
	{
		group := []view.View{}

		spacer := NewSpacer()
		l.Add(spacer, nil)

		switchView := view.NewSwitch()
		switchView.Value = v.app.Wifi.Enabled()
		switchView.OnSubmit = func(value bool) {
			v.app.Wifi.SetEnabled(!v.app.Wifi.Enabled())
		}

		cell1 := NewBasicCell()
		cell1.Title = "Wi-Fi"
		cell1.AccessoryView = switchView
		group = append(group, cell1)

		if v.app.Wifi.CurrentSSID() != "" && v.app.Wifi.Enabled() {
			cell2 := NewBasicCell()
			cell2.Title = v.app.Wifi.CurrentSSID()
			group = append(group, cell2)
		}

		for _, i := range AddSeparators(group, 60) {
			l.Add(i, nil)
		}
	}

	if v.app.Wifi.Enabled() {
		{
			group := []view.View{}

			spacer := NewSpacerHeader()
			spacer.Title = "Choose a Network..."
			l.Add(spacer, nil)

			for _, i := range v.app.Wifi.Networks() {
				network := i
				ssid := network.SSID()

				// Don't show the current network in this list.
				if ssid == v.app.Wifi.CurrentSSID() {
					continue
				}

				info := NewInfoButton()
				info.OnPress = func() {
					v.app.Stack.Push(NewWifiNetworkView(v.app, network))
				}

				cell := NewBasicCell()
				cell.Title = ssid
				cell.AccessoryView = info
				cell.OnTap = func() {
					v.app.Wifi.SetCurrentSSID(ssid)
				}
				group = append(group, cell)
			}

			cell1 := NewBasicCell()
			cell1.Title = "Other..."
			group = append(group, cell1)

			for _, i := range AddSeparators(group, 60) {
				l.Add(i, nil)
			}
		}
		{

			spacer := NewSpacer()
			l.Add(spacer, nil)

			switchView := view.NewSwitch()
			switchView.Value = v.app.Wifi.AskToJoin()
			switchView.OnSubmit = func(a bool) {
				v.app.Wifi.SetAskToJoin(a)
			}
			cell1 := NewBasicCell()
			cell1.Title = "Ask to Join Networks"
			cell1.AccessoryView = switchView

			for _, i := range AddSeparators([]view.View{cell1}, 60) {
				l.Add(i, nil)
			}
		}
		{
			spacer := NewSpacerDescription()
			spacer.Description = "Known networks will be joined automatically. If no known networks are available, you will have to manually join a network."
			l.Add(spacer, nil)
		}
	}

	scrollView := view.NewScrollView()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
		Options: []view.Option{
			&ios.StackBar{Title: "Wifi"},
			&android.StackBar{Title: "Wifi"},
		},
	}
}

type WifiNetworkView struct {
	view.Embed
	app     *App
	network *WifiNetwork
}

type a struct {
	app     *App
	network *WifiNetwork
}

func NewWifiNetworkView(app *App, network *WifiNetwork) *WifiNetworkView {
	return &WifiNetworkView{
		Embed:   view.NewEmbed(a{app: app, network: network}),
		app:     app,
		network: network,
	}
}

func (v *WifiNetworkView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.network)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.network)
	}
}

func (v *WifiNetworkView) Build(ctx view.Context) view.Model {
	props := v.network.Properties()

	l := &table.Layouter{}
	{

		spacer := NewSpacer()
		l.Add(spacer, nil)

		cell1 := NewBasicCell()
		cell1.Title = "Forget This Network"
		cell1.OnTap = func() {
			view.Alert("Forget Wi-Fi Network?", "Your iPhone will no longer join this Wi-Fi network.",
				&view.AlertButton{
					Title: "Cancel",
				},
				&view.AlertButton{
					Title: "Forget",
					OnPress: func() {
						v.app.Stack.Pop()
					},
				},
			)
		}

		for _, i := range AddSeparators([]view.View{cell1}, 60) {
			l.Add(i, nil)
		}
	}
	{

		spacer := NewSpacerHeader()
		spacer.Title = "IP Address"
		l.Add(spacer, nil)

		cell0 := NewSegmentCell()
		cell0.Titles = []string{"DHCP", "BootP", "Static"}
		cell0.Value = props.Kind
		cell0.OnValueChange = func(a int) {
			props := v.network.Properties()
			props.Kind = a
			v.network.SetProperties(props)
		}

		cell1 := NewBasicCell()
		cell1.Title = "IP Address"
		cell1.Subtitle = props.IPAddress

		cell2 := NewBasicCell()
		cell2.Title = "Subnet Mask"
		cell2.Subtitle = props.SubnetMask

		cell3 := NewBasicCell()
		cell3.Title = "Router"
		cell3.Subtitle = props.Router

		cell4 := NewBasicCell()
		cell4.Title = "DNS"
		cell4.Subtitle = props.DNS

		cell5 := NewBasicCell()
		cell5.Title = "Client ID"
		cell5.Subtitle = props.ClientID

		for _, i := range AddSeparators([]view.View{cell0, cell1, cell2, cell3, cell4, cell5}, 60) {
			l.Add(i, nil)
		}
	}
	{

		spacer := NewSpacer()
		l.Add(spacer, nil)

		cell1 := NewBasicCell()
		cell1.Title = "Renew Lease"
		cell1.OnTap = func() {
			view.Alert("Renewing Lease...", "")
		}

		for _, i := range AddSeparators([]view.View{cell1}, 60) {
			l.Add(i, nil)
		}
	}
	{

		spacer := NewSpacerHeader()
		spacer.Title = "HTTP Proxy"
		l.Add(spacer, nil)

		cell1 := NewSegmentCell()
		cell1.Titles = []string{"Off", "Manual", "Auto"}
		cell1.Value = props.Proxy
		cell1.OnValueChange = func(a int) {
			props := v.network.Properties()
			props.Proxy = a
			v.network.SetProperties(props)
		}

		for _, i := range AddSeparators([]view.View{cell1}, 60) {
			l.Add(i, nil)
		}
	}
	{

		spacer := NewSpacer()
		l.Add(spacer, nil)

		cell1 := NewBasicCell()
		cell1.Title = "Manage This Network"

		for _, i := range AddSeparators([]view.View{cell1}, 60) {
			l.Add(i, nil)
		}
	}
	spacer := NewSpacer()
	l.Add(spacer, nil)

	scrollView := view.NewScrollView()
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
		Options: []view.Option{
			&ios.StackBar{Title: v.network.SSID()},
			&android.StackBar{Title: v.network.SSID()},
		},
	}
}

type SegmentCell struct {
	view.Embed
	Titles        []string
	Value         int
	OnValueChange func(value int)
}

func NewSegmentCell() *SegmentCell {
	return &SegmentCell{}
}

func (v *SegmentCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(44)
		s.WidthEqual(l.MinGuide().Width())
	})

	segment := ios.NewSegmentView()
	segment.Titles = v.Titles
	segment.Value = v.Value
	segment.OnChange = func(a int) {
		if v.OnValueChange != nil {
			v.OnValueChange(a)
		}
	}
	l.Add(segment, func(s *constraint.Solver) {
		s.HeightLess(l.Height())
		s.RightEqual(l.Right().Add(-15))
		s.LeftEqual(l.Left().Add(15))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

type InfoButton struct {
	view.Embed
	OnPress    func()
	PaintStyle *paint.Style
}

func NewInfoButton() *InfoButton {
	return &InfoButton{}
}

func (v *InfoButton) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(44)
		s.Height(44)
	})

	img := view.NewImageView()
	img.Image = application.MustLoadImage("info")
	l.Add(img, func(s *constraint.Solver) {
		s.Width(22)
		s.Height(22)
		s.RightEqual(l.Right())
	})

	button := &pointer.ButtonGesture{
		OnEvent: func(e *pointer.ButtonEvent) {
			if e.Kind == pointer.EventKindRecognized && v.OnPress != nil {
				v.OnPress()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  v.PaintStyle,
		Options: []view.Option{
			pointer.GestureList{button},
		},
	}
}

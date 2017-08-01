package settings

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/alert"
	"gomatcha.io/matcha/view/imageview"
	"gomatcha.io/matcha/view/scrollview"
	"gomatcha.io/matcha/view/segmentview"
	"gomatcha.io/matcha/view/stackview"
	"gomatcha.io/matcha/view/switchview"
)

type WifiView struct {
	view.Embed
	app *App
}

func NewWifiView(ctx *view.Context, key string, app *App) *WifiView {
	if v, ok := ctx.Prev(key).(*WifiView); ok {
		return v
	}
	return &WifiView{
		Embed: ctx.NewEmbed(key),
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

func (v *WifiView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}
	{
		ctx := ctx.WithPrefix("1")
		group := []view.View{}

		spacer := NewSpacer(ctx, "spacer")
		l.Add(spacer, nil)

		switchView := switchview.New(ctx, "switch")
		switchView.Value = v.app.Wifi.Enabled()
		switchView.OnValueChange = func(value bool) {
			v.app.Wifi.SetEnabled(!v.app.Wifi.Enabled())
		}

		cell1 := NewBasicCell(ctx, "wifi")
		cell1.Title = "Wi-Fi"
		cell1.AccessoryView = switchView
		group = append(group, cell1)

		if v.app.Wifi.CurrentSSID() != "" && v.app.Wifi.Enabled() {
			cell2 := NewBasicCell(ctx, "current")
			cell2.Title = v.app.Wifi.CurrentSSID()
			group = append(group, cell2)
		}

		for _, i := range AddSeparators(ctx, group) {
			l.Add(i, nil)
		}
	}

	if v.app.Wifi.Enabled() {
		{
			ctx := ctx.WithPrefix("2")
			group := []view.View{}

			spacer := NewSpacerHeader(ctx, "spacer")
			spacer.Title = "Choose a Network..."
			l.Add(spacer, nil)

			for _, i := range v.app.Wifi.Networks() {
				network := i
				ssid := network.SSID()

				// Don't show the current network in this list.
				if ssid == v.app.Wifi.CurrentSSID() {
					continue
				}

				info := NewInfoButton(ctx, "networkbutton"+ssid)
				info.OnPress = func() {
					v.app.Stack.Push(NewWifiNetworkView(nil, "", v.app, network))
				}

				cell := NewBasicCell(ctx, "network"+ssid)
				cell.Title = ssid
				cell.AccessoryView = info
				cell.OnTap = func() {
					v.app.Wifi.SetCurrentSSID(ssid)
				}
				group = append(group, cell)
			}

			cell1 := NewBasicCell(ctx, "other")
			cell1.Title = "Other..."
			group = append(group, cell1)

			for _, i := range AddSeparators(ctx, group) {
				l.Add(i, nil)
			}
		}
		{
			ctx := ctx.WithPrefix("3")

			spacer := NewSpacer(ctx, "spacer")
			l.Add(spacer, nil)

			switchView := switchview.New(ctx, "switch")
			switchView.Value = v.app.Wifi.AskToJoin()
			switchView.OnValueChange = func(a bool) {
				v.app.Wifi.SetAskToJoin(a)
			}
			cell1 := NewBasicCell(ctx, "join")
			cell1.Title = "Ask to Join Networks"
			cell1.AccessoryView = switchView

			for _, i := range AddSeparators(ctx, []view.View{cell1}) {
				l.Add(i, nil)
			}
		}
		{
			spacer := NewSpacerDescription(ctx, "spacerDescr")
			spacer.Description = "Known networks will be joined automatically. If no known networks are available, you will have to manually join a network."
			l.Add(spacer, nil)
		}
	}

	scrollView := scrollview.New(ctx, "scroll")
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

func (v *WifiView) StackBar(ctx *view.Context) *stackview.Bar {
	return &stackview.Bar{Title: "Wi-Fi"}
}

type WifiNetworkView struct {
	view.Embed
	app     *App
	network *WifiNetwork
}

func NewWifiNetworkView(ctx *view.Context, key string, app *App, network *WifiNetwork) *WifiNetworkView {
	if v, ok := ctx.Prev(key).(*WifiNetworkView); ok {
		return v
	}
	return &WifiNetworkView{
		Embed:   ctx.NewEmbed(key),
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

func (v *WifiNetworkView) Build(ctx *view.Context) view.Model {
	props := v.network.Properties()

	l := &table.Layouter{}
	{
		ctx := ctx.WithPrefix("1")

		spacer := NewSpacer(ctx, "spacer")
		l.Add(spacer, nil)

		cell1 := NewBasicCell(ctx, "forget")
		cell1.Title = "Forget This Network"
		cell1.OnTap = func() {
			alert.Alert("Forget Wi-Fi Network?", "Your iPhone will no longer join this Wi-Fi network.",
				&alert.Button{
					Title: "Cancel",
					Style: alert.ButtonStyleCancel,
				},
				&alert.Button{
					Title: "Forget",
					OnPress: func() {
						v.app.Stack.Pop()
					},
				},
			)
		}

		for _, i := range AddSeparators(ctx, []view.View{cell1}) {
			l.Add(i, nil)
		}
	}
	{
		ctx := ctx.WithPrefix("2")

		spacer := NewSpacerHeader(ctx, "spacer")
		spacer.Title = "IP Address"
		l.Add(spacer, nil)

		cell0 := NewSegmentCell(ctx, "segment")
		cell0.Titles = []string{"DHCP", "BootP", "Static"}
		cell0.Value = props.Kind
		cell0.OnValueChange = func(a int) {
			props := v.network.Properties()
			props.Kind = a
			v.network.SetProperties(props)
		}

		cell1 := NewBasicCell(ctx, "ip")
		cell1.Title = "IP Address"
		cell1.Subtitle = props.IPAddress

		cell2 := NewBasicCell(ctx, "subnet")
		cell2.Title = "Subnet Mask"
		cell2.Subtitle = props.SubnetMask

		cell3 := NewBasicCell(ctx, "router")
		cell3.Title = "Router"
		cell3.Subtitle = props.Router

		cell4 := NewBasicCell(ctx, "dns")
		cell4.Title = "DNS"
		cell4.Subtitle = props.DNS

		cell5 := NewBasicCell(ctx, "clientid")
		cell5.Title = "Client ID"
		cell5.Subtitle = props.ClientID

		for _, i := range AddSeparators(ctx, []view.View{cell0, cell1, cell2, cell3, cell4, cell5}) {
			l.Add(i, nil)
		}
	}
	{
		ctx := ctx.WithPrefix("3")

		spacer := NewSpacer(ctx, "spacer")
		l.Add(spacer, nil)

		cell1 := NewBasicCell(ctx, "renew")
		cell1.Title = "Renew Lease"
		cell1.OnTap = func() {
			alert.Alert("Renewing Lease...", "")
		}

		for _, i := range AddSeparators(ctx, []view.View{cell1}) {
			l.Add(i, nil)
		}
	}
	{
		ctx := ctx.WithPrefix("4")

		spacer := NewSpacerHeader(ctx, "spacer")
		spacer.Title = "HTTP Proxy"
		l.Add(spacer, nil)

		cell1 := NewSegmentCell(ctx, "segment")
		cell1.Titles = []string{"Off", "Manual", "Auto"}
		cell1.Value = props.Proxy
		cell1.OnValueChange = func(a int) {
			props := v.network.Properties()
			props.Proxy = a
			v.network.SetProperties(props)
		}

		for _, i := range AddSeparators(ctx, []view.View{cell1}) {
			l.Add(i, nil)
		}
	}
	{
		ctx := ctx.WithPrefix("5")

		spacer := NewSpacer(ctx, "spacer")
		l.Add(spacer, nil)

		cell1 := NewBasicCell(ctx, "manage")
		cell1.Title = "Manage This Network"

		for _, i := range AddSeparators(ctx, []view.View{cell1}) {
			l.Add(i, nil)
		}
	}
	spacer := NewSpacer(ctx, "spacer")
	l.Add(spacer, nil)

	scrollView := scrollview.New(ctx, "scroll")
	scrollView.ContentChildren = l.Views()
	scrollView.ContentLayouter = l

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

func (v *WifiNetworkView) StackBar(*view.Context) *stackview.Bar {
	return &stackview.Bar{
		Title: v.network.SSID(),
	}
}

type SegmentCell struct {
	view.Embed
	Titles        []string
	Value         int
	OnValueChange func(value int)
}

func NewSegmentCell(ctx *view.Context, key string) *SegmentCell {
	if v, ok := ctx.Prev(key).(*SegmentCell); ok {
		return v
	}
	return &SegmentCell{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *SegmentCell) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(44)
		s.WidthEqual(l.MinGuide().Width())
	})

	segment := segmentview.New(ctx, "segment")
	segment.Titles = v.Titles
	segment.Value = v.Value
	segment.OnValueChange = func(a int) {
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

func NewInfoButton(ctx *view.Context, key string) *InfoButton {
	if v, ok := ctx.Prev(key).(*InfoButton); ok {
		return v
	}
	return &InfoButton{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *InfoButton) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Width(44)
		s.Height(44)
	})

	img := imageview.New(ctx, "image")
	img.Image = app.MustLoadImage("Info")
	l.Add(img, func(s *constraint.Solver) {
		s.Width(22)
		s.Height(22)
		s.RightEqual(l.Right())
	})

	button := &touch.ButtonRecognizer{
		OnTouch: func(e *touch.ButtonEvent) {
			if e.Kind == touch.EventKindRecognized && v.OnPress != nil {
				v.OnPress()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  v.PaintStyle,
		Options: []view.Option{
			touch.RecognizerList{button},
		},
	}
}

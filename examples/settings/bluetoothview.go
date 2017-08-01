package settings

import (
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/scrollview"
	"gomatcha.io/matcha/view/switchview"
)

type BluetoothView struct {
	view.Embed
	app *App
}

func NewBluetoothView(ctx *view.Context, key string, app *App) *BluetoothView {
	if v, ok := ctx.Prev(key).(*BluetoothView); ok {
		return v
	}
	v := &BluetoothView{
		Embed: ctx.NewEmbed(key),
		app:   app,
	}
	return v
}

func (v *BluetoothView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.app.Bluetooth)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.app.Bluetooth)
	}
}

func (v *BluetoothView) Build(ctx *view.Context) view.Model {
	l := &table.Layouter{}
	{
		ctx := ctx.WithPrefix("1")
		group := []view.View{}

		spacer := NewSpacer(ctx, "spacer")
		l.Add(spacer, nil)

		switchView := switchview.New(ctx, "switch")
		switchView.Value = v.app.Bluetooth.Enabled()
		switchView.OnValueChange = func(value bool) {
			v.app.Bluetooth.SetEnabled(!v.app.Bluetooth.Enabled())
		}

		cell1 := NewBasicCell(ctx, "wifi")
		cell1.Title = "Bluetooth"
		cell1.AccessoryView = switchView
		group = append(group, cell1)

		for _, i := range AddSeparators(ctx, group) {
			l.Add(i, nil)
		}
	}
	if v.app.Bluetooth.Enabled() {
		ctx := ctx.WithPrefix("2")
		group := []view.View{}

		spacer := NewSpacerHeader(ctx, "spacer")
		spacer.Title = "My Devices"
		l.Add(spacer, nil)

		for _, i := range v.app.Bluetooth.Devices() {
			device := i
			cell := NewBasicCell(ctx, "network"+device.SSID())
			cell.Title = device.SSID()
			if device.Connected() {
				cell.Subtitle = "Connected"
			} else {
				cell.Subtitle = "Not Connected"
			}
			cell.OnTap = func() {
				device.SetConnected(!device.Connected())
				v.Signal()
			}
			group = append(group, cell)
		}

		for _, i := range AddSeparators(ctx, group) {
			l.Add(i, nil)
		}
	}

	scrollView := scrollview.New(ctx, "b")
	scrollView.ContentLayouter = l
	scrollView.ContentChildren = l.Views()

	return view.Model{
		Children: []view.View{scrollView},
		Painter:  &paint.Style{BackgroundColor: backgroundColor},
	}
}

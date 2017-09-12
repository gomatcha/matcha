// Package alert implements basic alerts.
//
//  view.Alert("Title", "Message") // Has an OK button by default.
//  view.Alert("Title", "Message", &Button{
//      Title:"Cancel",
//      OnPress: func() {
//          // Do something
//      }
//  })
package view

import (
	"runtime"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/bridge"
	pbview "gomatcha.io/matcha/pb/view"
)

var alertMaxId int64
var alerts map[int64]*_alert

func init() {
	alerts = map[int64]*_alert{}
	bridge.RegisterFunc("gomatcha.io/matcha/view/alert onPress", func(id, idx int64) {
		alert, ok := alerts[id]
		if !ok {
			return
		}
		button := alert.Buttons[idx]
		if button.OnPress != nil {
			button.OnPress()
		}
	})
}

type _alert struct {
	Title   string
	Message string
	Buttons []*AlertButton
}

func (a *_alert) marshalProtobuf(id int64) *pbview.Alert {
	b := []*pbview.AlertButton{}
	for _, i := range a.Buttons {
		b = append(b, i.marshalProtobuf())
	}

	return &pbview.Alert{
		Id:      id,
		Title:   a.Title,
		Message: a.Message,
		Buttons: b,
	}
}

func (a *_alert) display() {
	alertMaxId += 1
	alerts[alertMaxId] = a

	data, err := proto.Marshal(a.marshalProtobuf(alertMaxId))
	if err != nil {
		return
	}
	if runtime.GOOS == "android" {
		bridge.Bridge("").Call("displayAlert", bridge.Bytes(data))
	} else if runtime.GOOS == "darwin" {
		bridge.Bridge("").Call("displayAlert:", bridge.Bytes(data))
	}
}

// Alert displays an alert with the given title, message and buttons. If no buttons are passed, a default OK button is created.
func Alert(title, message string, buttons ...*AlertButton) {
	if len(buttons) == 0 {
		buttons = []*AlertButton{&AlertButton{Title: "OK"}}
	}
	a := _alert{
		Title:   title,
		Message: message,
		Buttons: buttons,
	}
	a.display()
}

// AlertButton represents an alert button.
type AlertButton struct {
	Title   string
	OnPress func()
}

func (a *AlertButton) marshalProtobuf() *pbview.AlertButton {
	return &pbview.AlertButton{
		Title: a.Title,
	}
}

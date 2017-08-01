// Package alert implements basic alerts.
//
//  alert.Alert("Title", "Message") // Has an OK button by default.
//  alert.Alert("Title", "Message", &Button{
//      Title:"Cancel",
//      OnPress: func() {
//          // Do something
//      }
//  })
package alert

import (
	"github.com/gogo/protobuf/proto"
	"gomatcha.io/bridge"
	pbalert "gomatcha.io/matcha/pb/view/alert"
)

var maxId int64
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
	Buttons []*Button
}

func (a *_alert) marshalProtobuf(id int64) *pbalert.View {
	b := []*pbalert.Button{}
	for _, i := range a.Buttons {
		b = append(b, i.marshalProtobuf())
	}

	return &pbalert.View{
		Id:      id,
		Title:   a.Title,
		Message: a.Message,
		Buttons: b,
	}
}

func (a *_alert) display() {
	maxId += 1
	alerts[maxId] = a

	data, err := proto.Marshal(a.marshalProtobuf(maxId))
	if err != nil {
		return
	}
	bridge.Bridge().Call("displayAlert:", bridge.Bytes(data))
}

// Alert displays an alert with the given title, message and buttons. If no buttons are passed, a default OK button is created.
func Alert(title, message string, buttons ...*Button) {
	if len(buttons) == 0 {
		buttons = []*Button{&Button{Title: "OK"}}
	}
	a := _alert{
		Title:   title,
		Message: message,
		Buttons: buttons,
	}
	a.display()
}

// Button represents an alert button.
type Button struct {
	Title   string
	Style   ButtonStyle
	OnPress func()
}

func (a *Button) marshalProtobuf() *pbalert.Button {
	return &pbalert.Button{
		Title: a.Title,
		Style: pbalert.ButtonStyle(a.Style),
	}
}

// Alert button styles.
type ButtonStyle int

const (
	ButtonStyleDefault ButtonStyle = iota
	ButtonStyleCancel
	ButtonStyleDestructive
)

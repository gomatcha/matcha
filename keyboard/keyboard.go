// Package keyboard exposes access to displaying and hiding the keyboard.
//  input := textinput.New(ctx, "input")
//  input.Text = v.text
//  input.KeyboardType = keyboard.URLType
//  input.Responder = v.responder
//
//  button := ...
//  button.OnTap = func() {
//  	v.responder.Dismiss()
//  }
package keyboard

import (
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/proto/keyboard"
)

// Type defines the kind of keyboard.
type Type int

const (
	TextType Type = iota
	NumberType
	EmailType
	URLType
	PhoneType
	// DateTimeType
)

func (t Type) MarshalProtobuf() keyboard.Type {
	return keyboard.Type(t)
}

// // ReturnType defines the keyboard return key style
// type ReturnType int

// const (
// 	DefaultReturnType ReturnType = iota
// 	GoReturnType
// 	GoogleReturnType
// 	JoinReturnType
// 	NextReturnType
// 	RouteReturnType
// 	SearchReturnType
// 	SendReturnType
// 	YahooReturnType
// 	DoneReturnType
// 	EmergencyCallReturnType
// 	ContinueReturnType
// )

// func (t ReturnType) MarshalProtobuf() keyboard.ReturnType {
// 	return keyboard.ReturnType(t)
// }

// Responder is a model object that represents the keyboard's state. To use Responder it must be attached to a textinput.View.
type Responder struct {
	visible bool
	value   comm.Relay
}

// Show displays the keyboard.
func (g *Responder) Show() {
	if !g.visible {
		g.visible = true
		g.value.Signal()
	}
}

// Dismiss hides any displayed keyboards.
func (g *Responder) Dismiss() {
	if g.visible {
		g.visible = false
		g.value.Signal()
	}
}

// Visible returns true if the keyboard is visible.
func (g *Responder) Visible() bool {
	return g.visible
}

// Notify implements comm.Notifier.
func (g *Responder) Notify(f func()) comm.Id {
	return g.value.Notify(f)
}

// Unnotify implements comm.Notifier.
func (g *Responder) Unnotify(id comm.Id) {
	g.value.Unnotify(id)
}

// func (g *Responder) Next() {
// }

// func (g *Responder) Prev() {
// }

// type key struct{}

// var Key = key{}

// type Middleware struct {
// 	radix *radix.Radix
// }

// func NewMiddleware() *Middleware {
// 	return &Middleware{radix: radix.NewRadix()}
// }

// func (m *Middleware) Build(ctx view.Context, next *view.Model) {
// 	responder, ok := next.Values[Key].(*Responder)
// 	path := []int64{}
// 	for _, i := range ctx.Path() {
// 		path = append(path, int64(i))
// 	}

// 	if ok {
// 		n := m.radix.Insert(path)
// 		n.Value = responder
// 	} else {
// 		m.radix.Delete(path)
// 	}
// }

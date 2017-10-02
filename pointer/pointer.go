/*
Package pointer implements gesture recognizers.

Create the pointer recognizer in the Build function.

 func (v *MyView) Build(ctx view.Context) view.Model {
 	tap := &pointer.TapGesture{
 		Count: 1,
 		OnEvent: func(e *pointer.TapEvent) {
 			if e.Kind == pointer.EventKindRecognized {
 				// Respond to pointer events. This callback occurs on main thread.
 				fmt.Println("view tapped")
			}
 		},
 	}
 	...

Attach the recognizer to the view.

	...
 	return view.Model{
 		Options: []view.Option{
 			pointer.GestureList{tap},
 		},
 	}
 }
*/
package pointer

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gomatcha.io/matcha/layout"
	pbtouch "gomatcha.io/matcha/proto/pointer"
)

var maxFuncId int64 = 0

// NewFuncId generates a new func identifier for serialization.
func newFuncId() int64 {
	return atomic.AddInt64(&maxFuncId, 1)
}

type GestureList []Gesture

func (r GestureList) OptionKey() string {
	return "gomatcha.io/matcha/touch"
}

type Gesture interface {
	Build() Model
	TouchKey() int64
}

type Model struct {
	NativeViewName  string
	NativeViewState proto.Message
	NativeFuncs     map[string]interface{}
}

// EventKind are the possible recognizer states
//
// Discrete gestures:
//  EventKindPossible -> EventKindFailed
//  EventKindPossible -> EventKindRecognized
// Continuous gestures:
//  EventKindPossible -> EventKindFailed
//  EventKindPossible -> EventKindRecognized
//  EventKindPossible -> EventKindChanged(optionally) -> EventKindFailed
//  EventKindPossible -> EventKindChanged(optionally) -> EventKindRecognized
type EventKind int

const (
	// Finger is down, but before gesture has been recognized.
	EventKindPossible EventKind = iota
	// After the continuous gesture has been recognized, while the finger is still down. Only for continuous recognizers.
	EventKindChanged
	// Gesture recognition failed or cancelled.
	EventKindFailed
	// Gesture recognition succeded.
	EventKindRecognized
)

// TapEvent is emitted by TapGesture, representing its current state.
type TapEvent struct {
	Kind      EventKind
	Timestamp time.Time
	Position  layout.Point
}

func (e *TapEvent) unmarshalProtobuf(ev *pbtouch.TapEvent) error {
	t, _ := ptypes.Timestamp(ev.Timestamp)
	e.Timestamp = t
	e.Kind = EventKind(ev.Kind)
	e.Position.UnmarshalProtobuf(ev.Position)
	return nil
}

// PressRecognizer is a discrete recognizer that detects a number of taps.
type TapGesture struct {
	Key     int64
	Count   int
	OnEvent func(*TapEvent)
}

func (r *TapGesture) TouchKey() int64 {
	return r.Key
}

func (r *TapGesture) Build() Model {
	funcId := newFuncId()
	f := func(data []byte) {
		pbevent := &pbtouch.TapEvent{}
		err := proto.Unmarshal(data, pbevent)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		event := &TapEvent{}
		if err := event.unmarshalProtobuf(pbevent); err != nil {
			fmt.Println("error", err)
			return
		}

		if r.OnEvent != nil {
			r.OnEvent(event)
		}
	}

	return Model{
		NativeViewName: "",
		NativeViewState: &pbtouch.TapRecognizer{
			Count:   int64(r.Count),
			OnEvent: funcId,
		},
		NativeFuncs: map[string]interface{}{
			fmt.Sprintf("gomatcha.io/matcha/touch %v", funcId): f,
		},
	}
}

// PressEvent is emitted by PressRecognizer, representing its current state.
type PressEvent struct {
	Kind      EventKind // TODO(KD): Does this work?
	Timestamp time.Time
	Position  layout.Point
	Duration  time.Duration
}

func (e *PressEvent) unmarshalProtobuf(ev *pbtouch.PressEvent) error {
	d, err := ptypes.Duration(ev.Duration)
	if err != nil {
		return err
	}
	t, err := ptypes.Timestamp(ev.Timestamp)
	if err != nil {
		return err
	}
	e.Kind = EventKind(ev.Kind)
	e.Timestamp = t
	e.Position.UnmarshalProtobuf(ev.Position)
	e.Duration = d
	return nil
}

// PressGesture is a continuous recognizer that detects single presses with a given duration.
type PressGesture struct {
	Key         int64
	MinDuration time.Duration
	OnEvent     func(e *PressEvent)
}

func (r *PressGesture) TouchKey() int64 {
	return r.Key
}

func (r *PressGesture) Build() Model {
	funcId := newFuncId()
	f := func(data []byte) {
		event := &PressEvent{}
		pbevent := &pbtouch.PressEvent{}
		err := proto.Unmarshal(data, pbevent)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		if err := event.unmarshalProtobuf(pbevent); err != nil {
			fmt.Println("error", err)
			return
		}
		if r.OnEvent != nil {
			r.OnEvent(event)
		}
	}

	return Model{
		NativeViewName: "",
		NativeViewState: &pbtouch.PressRecognizer{
			MinDuration: ptypes.DurationProto(r.MinDuration),
			OnEvent:     funcId,
		},
		NativeFuncs: map[string]interface{}{
			fmt.Sprintf("gomatcha.io/matcha/touch %v", funcId): f,
		},
	}
}

// ButtonEvent is emitted by ButtonRecognizer, representing its current state.
type ButtonEvent struct {
	Timestamp time.Time
	Inside    bool
	Kind      EventKind
}

func (e *ButtonEvent) unmarshalProtobuf(ev *pbtouch.ButtonEvent) error {
	t, err := ptypes.Timestamp(ev.Timestamp)
	if err != nil {
		return err
	}
	e.Timestamp = t
	e.Inside = ev.Inside
	e.Kind = EventKind(ev.Kind)
	return nil
}

// ButtonGesture is a discrete recognizer that mimics the behavior of a button. The recognizer will fail if the touch ends outside of the view's bounds.
type ButtonGesture struct {
	Key           int64
	OnEvent       func(e *ButtonEvent)
	IgnoresScroll bool
}

func (r *ButtonGesture) TouchKey() int64 {
	return r.Key
}

func (r *ButtonGesture) Build() Model {
	funcId := newFuncId()
	f := func(data []byte) {
		event := &ButtonEvent{}
		pbevent := &pbtouch.ButtonEvent{}
		err := proto.Unmarshal(data, pbevent)
		if err != nil {
			fmt.Println("error", err)
			return
		}

		if err := event.unmarshalProtobuf(pbevent); err != nil {
			fmt.Println("error", err)
			return
		}

		if r.OnEvent != nil {
			r.OnEvent(event)
		}
	}

	return Model{
		NativeViewName: "",
		NativeViewState: &pbtouch.ButtonRecognizer{
			OnEvent:       funcId,
			IgnoresScroll: r.IgnoresScroll,
		},
		NativeFuncs: map[string]interface{}{
			fmt.Sprintf("gomatcha.io/matcha/touch %v", funcId): f,
		},
	}
}

/*
Package constraint implements touch recognizers.

Create the touch recognizer in the Build function.

 func (v *MyView) Build(ctx *view.Context) view.Model {
 	tap := &touch.TapRecognizer{
 		Count: 1,
 		OnTouch: func(e *touch.TapEvent) {
 			// Respond to touch events. This callback occurs on main thread.
 			fmt.Println("view touched")
 		},
 	}
 	...

Attach the recognizer to the view.

	...
 	return view.Model{
 		Options: []view.Option{
 			touch.RecognizerList{tap},
 		},
 	}
 }
*/
package touch

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gomatcha.io/matcha/layout"
	pbtouch "gomatcha.io/matcha/pb/touch"
)

var maxFuncId int64 = 0

// NewFuncId generates a new func identifier for serialization.
func newFuncId() int64 {
	return atomic.AddInt64(&maxFuncId, 1)
}

type RecognizerList []Recognizer

func (r RecognizerList) OptionKey() string {
	return "gomatcha.io/matcha/touch"
}

type Recognizer interface {
	marshalProtobuf() (proto.Message, map[string]interface{})
	equal(Recognizer) bool
}

// TapEvent is emitted by TapRecognizer, representing its current state.
type TapEvent struct {
	// Kind      EventKind // TODO(KD):

	Timestamp time.Time
	Position  layout.Point
}

func (e *TapEvent) unmarshalProtobuf(ev *pbtouch.TapEvent) error {
	t, _ := ptypes.Timestamp(ev.Timestamp)
	e.Timestamp = t
	// e.Kind = EventKind(ev.Kind)
	e.Position.UnmarshalProtobuf(ev.Position)
	return nil
}

// PressRecognizer is a discrete recognizer that detects a number of taps.
type TapRecognizer struct {
	Count   int
	OnTouch func(*TapEvent)
}

func (r *TapRecognizer) equal(a Recognizer) bool {
	b, ok := a.(*TapRecognizer)
	if !ok {
		return false
	}
	return r.Count == b.Count
}

func (r *TapRecognizer) marshalProtobuf() (proto.Message, map[string]interface{}) {
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

		if r.OnTouch != nil {
			r.OnTouch(event)
		}
	}

	return &pbtouch.TapRecognizer{
			Count:          int64(r.Count),
			RecognizedFunc: funcId,
		}, map[string]interface{}{
			strconv.Itoa(int(funcId)): f,
		}
}

// EventKind are the possible recognizer states
//
// Discrete gestures:
//  EventKindPossible -> EventKindFailed
//  EventKindPossible -> EventKindRecognized
// Continuous gestures:
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

// PressRecognizer is a continuous recognizer that detects single presses with a given duration.
type PressRecognizer struct {
	MinDuration time.Duration
	OnTouch     func(e *PressEvent)
}

func (r *PressRecognizer) equal(a Recognizer) bool {
	b, ok := a.(*PressRecognizer)
	if !ok {
		return false
	}
	return r.MinDuration == b.MinDuration
}

func (r *PressRecognizer) marshalProtobuf() (proto.Message, map[string]interface{}) {
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
		if r.OnTouch != nil {
			r.OnTouch(event)
		}
	}

	return &pbtouch.PressRecognizer{
			MinDuration: ptypes.DurationProto(r.MinDuration),
			FuncId:      funcId,
		}, map[string]interface{}{
			strconv.Itoa(int(funcId)): f,
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

// ButtonRecognizer is a discrete recognizer that mimics the behavior of a button. The recognizer will fail if the touch ends outside of the view's bounds.
type ButtonRecognizer struct {
	OnTouch       func(e *ButtonEvent)
	IgnoresScroll bool
}

func (r *ButtonRecognizer) equal(a Recognizer) bool {
	_, ok := a.(*ButtonRecognizer)
	if !ok {
		return false
	}
	return true
}

func (r *ButtonRecognizer) marshalProtobuf() (proto.Message, map[string]interface{}) {
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

		if r.OnTouch != nil {
			r.OnTouch(event)
		}
	}

	return &pbtouch.ButtonRecognizer{
			OnEvent:       funcId,
			IgnoresScroll: r.IgnoresScroll,
		}, map[string]interface{}{
			strconv.Itoa(int(funcId)): f,
		}
}

package pointer

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"gomatcha.io/matcha/layout"
	prototouch "gomatcha.io/matcha/proto/pointer"
)

// gestureState are the possible recognizer states
//
// Discrete gestures:
//  gestureStatePossible -> gestureStateFailed
//  gestureStatePossible -> gestureStateRecognized
// Continuous gestures:
//  gestureStatePossible -> gestureStateChanged(optionally) -> gestureStateFailed
//  gestureStatePossible -> gestureStateChanged(optionally) -> gestureStateRecognized
type gestureState int

const (
	// Finger is down, but before gesture has been recognized.
	gestureStatePossible gestureState = iota
	// After the continuous gesture has been recognized, while the finger is still down. Only for continuous recognizers.
	gestureStateChanged
	// Gesture recognition failed or cancelled.
	gestureStateFailed
	// Gesture recognition succeded.
	gestureStateRecognized
)

type phase int

const (
	phaseBegan phase = iota
	phaseMoved
	phaseEnded
	phaseCancelled
)

type event struct {
	Timestamp time.Time
	Location  layout.Point
	Phase     phase
	// Pointers []pointer Multitouch....
}

func (e *event) unmarshalProtobuf(ev *prototouch.Event) error {
	t, _ := ptypes.Timestamp(ev.Timestamp)
	e.Timestamp = t
	e.Phase = phase(ev.Phase)
	e.Location.UnmarshalProtobuf(ev.Location)
	return nil
}

type gesture2 interface {
	onEvent(e *event) gestureState
	reset()
	gestureKey() interface{}

	// Future...
	// priority() int
}

type TapGesture2 struct {
	Count       int
	OnRecognize func(e *TapEvent2)
}

func (r *TapGesture2) OptionKey() string {
	return "gomatcha.io/matcha/touch Gesture"
}

func (r *TapGesture2) gestureKey() interface{} {
	return r.Count
}

func (r *TapGesture2) onEvent(e *event) gestureState {
	if e.Phase == phaseBegan {
		return gestureStatePossible
	} else if e.Phase == phaseEnded {
		return gestureStateRecognized
	} else if e.Phase == phaseCancelled {
		return gestureStateFailed
	}
	return gestureStatePossible
}

func (r *TapGesture2) reset() {
	fmt.Println("reset")
}

type TapEvent2 struct {
	Timestamp time.Time
	Location  layout.Point
}

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
	phaseNone
)

type event struct {
	Timestamp time.Time
	Location  layout.Point
	Phase     phase
	ViewSize  layout.Point
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
	onEvent(gestureState, *event) gestureState
	reset()
	gestureKey() interface{}

	// Future...
	// priority() int
}

const tapMaxDuration = time.Duration(0.5 * float64(time.Second))
const tapMaxDurationBetween = time.Duration(0.5 * float64(time.Second))
const tapMaxDistance = 25
const tapMaxDistanceBetween = 50

type TapGesture struct {
	Count       int
	OnRecognize func(e *TapEvent)

	state          gestureState
	count          int
	startTimestamp []time.Time
	startLocation  []layout.Point
}

func (g *TapGesture) OptionKey() string {
	return "gomatcha.io/matcha/touch Gesture"
}

func (g *TapGesture) gestureKey() interface{} {
	return g.Count
}

func (g *TapGesture) onEvent(s gestureState, e *event) gestureState {
	if s != g.state {
		g.state = gestureStateFailed
		return gestureStateFailed
	}

	switch e.Phase {
	case phaseBegan:
		if g.count > 0 {
			timestamp := g.startTimestamp[len(g.startTimestamp)-1]
			location := g.startLocation[len(g.startLocation)-1]
			if e.Timestamp.Sub(timestamp) > tapMaxDurationBetween || e.Location.Sub(location).Norm() > tapMaxDistanceBetween {
				g.state = gestureStateFailed
				return gestureStateFailed
			}
		}

		g.startTimestamp = append(g.startTimestamp, e.Timestamp)
		g.startLocation = append(g.startLocation, e.Location)
		g.state = gestureStatePossible
		return gestureStatePossible
	case phaseMoved:
		timestamp := g.startTimestamp[len(g.startTimestamp)-1]
		location := g.startLocation[len(g.startLocation)-1]
		if e.Timestamp.Sub(timestamp) > tapMaxDuration || e.Location.Sub(location).Norm() > tapMaxDistance {
			g.state = gestureStateFailed
			return gestureStateFailed
		}

		g.state = gestureStatePossible
		return gestureStatePossible
	case phaseEnded:
		timestamp := g.startTimestamp[len(g.startTimestamp)-1]
		location := g.startLocation[len(g.startLocation)-1]
		if e.Timestamp.Sub(timestamp) > tapMaxDuration || e.Location.Sub(location).Norm() > tapMaxDistance {
			g.state = gestureStateFailed
			return gestureStateFailed
		}

		if g.count == g.Count-1 {
			if g.OnRecognize != nil {
				tapEvent := &TapEvent{
					Timestamp: e.Timestamp,
					Location:  e.Location,
				}
				g.OnRecognize(tapEvent)
			}
			g.state = gestureStateRecognized
			return gestureStateRecognized
		} else if g.count < g.Count-1 {
			g.count += 1
			g.state = gestureStatePossible
			return gestureStatePossible
		}
		g.state = gestureStateFailed
		return gestureStateFailed
	case phaseNone:
		if g.count > 0 {
			timestamp := g.startTimestamp[len(g.startTimestamp)-1]
			if e.Timestamp.Sub(timestamp) > tapMaxDurationBetween {
				g.state = gestureStateFailed
				return gestureStateFailed
			}
		}

		g.state = gestureStatePossible
		return gestureStatePossible
	default:
		g.state = gestureStateFailed
		return gestureStateFailed
	}
}

func (g *TapGesture) reset() {
	g.state = gestureStatePossible
	g.startTimestamp = nil
	g.startLocation = nil
	g.count = 0
}

type TapEvent struct {
	Timestamp time.Time
	Location  layout.Point
}

type ButtonGesture2 struct {
	OnHighlight func(highlighted bool)
	OnRecognize func(e *ButtonEvent2)
	highlighted bool
}

func (g *ButtonGesture2) OptionKey() string {
	return "gomatcha.io/matcha/touch Gesture"
}

func (g *ButtonGesture2) gestureKey() interface{} {
	return nil
}

func (g *ButtonGesture2) onEvent(s gestureState, e *event) gestureState {
	switch e.Phase {
	case phaseBegan:
		highlighted := e.Location.X >= 0 && e.Location.Y >= 0 && e.Location.X <= e.ViewSize.X && e.Location.Y <= e.ViewSize.Y
		if g.OnHighlight != nil && g.highlighted != highlighted {
			g.OnHighlight(highlighted)
		}
		g.highlighted = highlighted

		return gestureStatePossible
	case phaseMoved:
		highlighted := e.Location.X >= 0 && e.Location.Y >= 0 && e.Location.X <= e.ViewSize.X && e.Location.Y <= e.ViewSize.Y
		if g.OnHighlight != nil && g.highlighted != highlighted {
			g.OnHighlight(highlighted)
		}
		g.highlighted = highlighted

		return gestureStatePossible
	case phaseEnded:
		highlighted := e.Location.X >= 0 && e.Location.Y >= 0 && e.Location.X <= e.ViewSize.X && e.Location.Y <= e.ViewSize.Y
		if g.OnHighlight != nil && g.highlighted != highlighted {
			g.OnHighlight(highlighted)
		}
		g.highlighted = highlighted

		if !highlighted {
			return gestureStateFailed
		}

		if g.OnRecognize != nil {
			buttonEvent := &ButtonEvent2{
				Timestamp: e.Timestamp,
				Location:  e.Location,
			}
			g.OnRecognize(buttonEvent)
		}
		return gestureStateRecognized
	case phaseNone:
		return gestureStateFailed
	default:
		return gestureStateFailed
	}
}

func (g *ButtonGesture2) reset() {
	if g.OnHighlight != nil && g.highlighted {
		g.OnHighlight(false)
	}
	g.highlighted = false
}

type ButtonEvent2 struct {
	Timestamp time.Time
	Location  layout.Point
}

type PressGesture struct {
	MinDuration time.Duration
	OnChange    func(e *PressEvent)
	OnFail      func(e *PressEvent)
	OnRecognize func(e *PressEvent)
}

type PressEvent struct {
	Timestamp time.Time
	Location  layout.Point
	Duration  time.Duration
}

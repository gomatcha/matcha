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
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	"gomatcha.io/matcha/layout"
	pbtouch "gomatcha.io/matcha/pb/touch"
	"gomatcha.io/matcha/view"
)

func init() {
	internal.RegisterMiddleware(func() interface{} { return &middleware{radix: radix.NewRadix()} })
}

type middleware struct {
	maxId int64
	radix *radix.Radix
}

var maxFuncId int64 = 0

// NewFuncId generates a new func identifier for serialization.
func newFuncId() int64 {
	return atomic.AddInt64(&maxFuncId, 1)
}

func (r *middleware) MarshalProtobuf() proto.Message {
	return nil
}

func (r *middleware) Build(ctx *view.Context, next *view.Model) {
	path := idSliceToIntSlice(ctx.Path())
	node := r.radix.At(path)
	var prevIds map[int64]Recognizer
	if node != nil {
		prevIds, _ = node.Value.(map[int64]Recognizer)
	}

	ids := map[int64]Recognizer{}

	var rs RecognizerList
	for _, i := range next.Options {
		rs, _ = i.(RecognizerList)
		if rs != nil {
			break
		}
	}

	// Diff prev and next recognizers
	for _, i := range rs {
		found := false
		for k, v := range prevIds {
			// Check that the id has not already been used.
			if _, ok := ids[k]; ok {
				continue
			}

			// Check that the recognizers are equal.
			if !i.equal(v) {
				continue
			}

			ids[k] = i
			found = true
		}

		// Generate a new id if we don't have a previous one.
		if !found {
			r.maxId += 1
			ids[r.maxId] = i
		}
	}

	if len(ids) == 0 {
		r.radix.Delete(path)
		return
	}

	// Add new list back to next.
	if node == nil {
		node = r.radix.Insert(path)
	}
	node.Value = ids

	// Serialize into protobuf.
	pbRecognizers := &pbtouch.RecognizerList{}
	allFuncs := map[string]interface{}{}
	for k, v := range ids {
		msg, funcs := v.marshalProtobuf(ctx)
		pbAny, err := ptypes.MarshalAny(msg)
		if err != nil {
			continue
		}

		pbRecognizer := &pbtouch.Recognizer{
			Id:         k,
			Recognizer: pbAny,
		}
		pbRecognizers.Recognizers = append(pbRecognizers.Recognizers, pbRecognizer)
		for k2, v2 := range funcs {
			allFuncs[k2] = v2
		}
	}

	if next.NativeValues == nil {
		next.NativeValues = map[string]proto.Message{}
	}
	next.NativeValues["gomatcha.io/matcha/touch"] = pbRecognizers

	if next.NativeFuncs == nil {
		next.NativeFuncs = map[string]interface{}{}
	}
	for k, v := range allFuncs {
		next.NativeFuncs[k] = v
	}
}

func (r *middleware) Key() string {
	return "gomatcha.io/matcha/touch"
}

type RecognizerList []Recognizer

func (r RecognizerList) OptionsKey() string {
	return "gomatcha.io/matcha/touch"
}

type Recognizer interface {
	marshalProtobuf(ctx *view.Context) (proto.Message, map[string]interface{})
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

func (r *TapRecognizer) marshalProtobuf(ctx *view.Context) (proto.Message, map[string]interface{}) {
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

func (r *PressRecognizer) marshalProtobuf(ctx *view.Context) (proto.Message, map[string]interface{}) {
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

func (r *ButtonRecognizer) marshalProtobuf(ctx *view.Context) (proto.Message, map[string]interface{}) {
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

func idSliceToIntSlice(ids []matcha.Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

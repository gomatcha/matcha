package pointer

import (
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	protopointer "gomatcha.io/matcha/proto/pointer"
	"gomatcha.io/matcha/view"
)

func init() {
	internal.RegisterMiddleware(func(root *internal.MiddlewareRoot) interface{} {
		m := &touchMiddleware{
			root:  root,
			radix: radix.NewRadix(),
		}
		t := internal.NewTicker(func() {
			m.tick()
		})
		m.ticker = t
		return m
	})
}

type key struct {
	gestureKey  interface{}
	reflectName string
}

type gestureWrapper struct {
	gesture gesture2
	state   gestureState
}

type gestureNode struct {
	middleware *touchMiddleware
	gestures   []*gestureWrapper
	getTicks   bool
	id         view.Id
}

func (n *gestureNode) validate() bool {
	possible := 0
	recognized := 0
	changed := 0
	for _, i := range n.gestures {
		switch i.state {
		case gestureStatePossible:
			possible += 1
		case gestureStateChanged:
			changed += 1
		case gestureStateRecognized:
			recognized += 1
		case gestureStateFailed:
		default:
			return false
		}
	}
	// There can be only one recognized gesture or changing gesture.
	if changed == 1 && recognized == 1 {
		return false
	}
	// If there is a recognized gesture or changing gesture, there can be no possible gestures.
	if possible > 0 && (changed == 1 || recognized == 1) {
		return false
	}
	return true
}

func (n *gestureNode) state() gestureState {
	possible := 0
	recognized := 0
	changed := 0
	failed := 0
	for _, i := range n.gestures {
		switch i.state {
		case gestureStatePossible:
			possible += 1
		case gestureStateChanged:
			changed += 1
		case gestureStateRecognized:
			recognized += 1
		case gestureStateFailed:
			failed += 1
		default:
			panic("Internal inconsistency")
		}
	}

	if recognized == 1 {
		return gestureStateRecognized
	} else if changed == 1 {
		return gestureStateChanged
	} else if possible > 0 {
		return gestureStatePossible
	} else if failed == len(n.gestures) {
		return gestureStateFailed
	} else {
		panic("Internal inconsistency")
	}
}

func (n *gestureNode) setPossible(idx int) {
	g := n.gestures[idx]
	if g.state != gestureStatePossible {
		panic("Internal inconsistency")
	}
	g.state = gestureStatePossible
}

func (n *gestureNode) setChanged(idx int) {
	g := n.gestures[idx]
	if g.state != gestureStatePossible && g.state != gestureStateChanged {
		panic("Internal inconsistency")
	}
	g.state = gestureStateChanged
}

func (n *gestureNode) setFailed(idx int) {
	g := n.gestures[idx]
	if g.state != gestureStatePossible && g.state != gestureStateChanged {
		fmt.Println("g.state", g.state)
		panic("Internal inconsistency")
	}
	g.gesture.reset()
	g.state = gestureStateFailed
}

func (n *gestureNode) setRecognized(ridx int) {
	g := n.gestures[ridx]
	if g.state != gestureStatePossible && g.state != gestureStateChanged {
		panic("Internal inconsistency")
	}
	g.gesture.reset()
	g.state = gestureStateRecognized

	// mark any possible gestures as failed
	for idx, i := range n.gestures {
		if idx != ridx && i.state == gestureStatePossible {
			i.gesture.reset()
			i.state = gestureStateFailed
		}
	}
}

func (n *gestureNode) markComplete() {
	for _, i := range n.gestures {
		if i.state == gestureStateRecognized {
			i.state = gestureStateFailed
		}
	}
}

func (n *gestureNode) markPossible() {
	for _, i := range n.gestures {
		if i.state == gestureStateFailed {
			i.state = gestureStatePossible
		} else {
			panic("Internal inconsistency")
		}
	}
}

func (n *gestureNode) onEvent(e *event) (bool, gestureState) {
	if !n.validate() {
		panic("Internal inconsistency")
	}

	if n.state() == gestureStateFailed && e.Phase == phaseBegan {
		n.markPossible()
	}

	for idx, i := range n.gestures {
		done := false

		switch i.state {
		case gestureStatePossible, gestureStateChanged:
			state := i.gesture.onEvent(i.state, e)
			switch state {
			case gestureStatePossible:
				n.setPossible(idx)
			case gestureStateChanged:
				n.setChanged(idx)
				done = true
			case gestureStateFailed:
				n.setFailed(idx)
			case gestureStateRecognized:
				n.setRecognized(idx)
				done = true
			default:
				panic("Internal inconsistency")
			}
		}

		if done {
			break
		}
	}

	state := n.state()

	// Reset if recognized
	switch state {
	case gestureStateRecognized:
		n.markComplete()
	}

	// If touch is up but the gesture is still possible, send a event on every screen update with phaseNone so the gesture can cancel.
	if (state == gestureStatePossible || state == gestureStateChanged) && (e.Phase == phaseEnded || e.Phase == phaseCancelled || e.Phase == phaseNone) {
		if !n.getTicks {
			n.getTicks = true
			n.middleware.tickNodes = append(n.middleware.tickNodes, n)
		}
	} else {
		if n.getTicks {
			n.getTicks = false

			tickNodes := make([]*gestureNode, 0, len(n.middleware.tickNodes))
			for _, i := range n.middleware.tickNodes {
				if i != n {
					tickNodes = append(tickNodes, i)
				}
			}
			n.middleware.tickNodes = tickNodes
		}
	}

	return false, state
}

func (n *gestureNode) reset() {
	for _, i := range n.gestures {
		switch i.state {
		case gestureStatePossible, gestureStateChanged:
			i.gesture.reset()
			i.state = gestureStateFailed
		case gestureStateRecognized:
			i.state = gestureStateFailed
		case gestureStateFailed:
		default:
			panic("Internal inconsistency")
		}
	}
}

type touchMiddleware struct {
	root      *internal.MiddlewareRoot
	maxId     int64
	radix     *radix.Radix
	ticker    *internal.Ticker
	tickNodes []*gestureNode // nodes that want a tick message
}

func (r *touchMiddleware) MarshalProtobuf() proto.Message {
	return nil
}

func (r *touchMiddleware) Build(ctx view.Context, model *view.Model) {
	// Get current gestures
	curr := []*gestureWrapper{}
	if model != nil {
		for _, i := range model.Options {
			g, _ := i.(gesture2)
			if i.OptionKey() == "gomatcha.io/matcha/touch Gesture" && g != nil {
				curr = append(curr, &gestureWrapper{gesture: g, state: gestureStateFailed})
			}
		}
	}

	// Delete the node if we have no gestures
	path := idSliceToIntSlice(ctx.Path())
	if len(curr) == 0 {
		r.radix.Delete(path)
		return
	}

	// Create the node, or get the previous node.
	node := r.radix.At(path)
	if node == nil {
		node = r.radix.Insert(path)
	}
	n, _ := node.Value.(*gestureNode)
	if n == nil {
		n = &gestureNode{
			middleware: r,
			id:         view.Id(path[len(path)-1]),
		}
		node.Value = n
	}

	// Diff prev and curr recognizers
	prev := n.gestures
	currKeys := make([]interface{}, len(curr))
	prevKeys := make([]interface{}, len(prev))
	for idx, i := range curr {
		currKeys[idx] = key{i.gesture.gestureKey(), internal.ReflectName(i.gesture)}
	}
	for idx, i := range prev {
		prevKeys[idx] = key{i.gesture.gestureKey(), internal.ReflectName(i.gesture)}
	}
	_, matching, _ := internal.Diff(currKeys, prevKeys)
	for k, v := range matching {
		// replace the current gesture with the previous gesture if they are matching.
		curr[k] = prev[v]
	}

	// Update node
	n.gestures = curr

	if model.NativeOptions == nil {
		model.NativeOptions = map[string][]byte{}
	}
	model.NativeOptions["gomatcha.io/matcha/pointer Gestures"] = []byte{}

	if model.NativeFuncs == nil {
		model.NativeFuncs = map[string]interface{}{}
	}
	model.NativeFuncs["gomatcha.io/matcha/pointer OnEvent"] = func(v []byte) (bool, int64) {
		protoEvent := &protopointer.Event{}
		if err := proto.Unmarshal(v, protoEvent); err != nil {
			fmt.Println("error", err)
			return false, 0
		}

		layoutGuide := r.root.LayoutGuide(int64(n.id))
		e := &event{}
		e.unmarshalProtobuf(protoEvent)
		e.ViewSize = layoutGuide.Frame.Size()

		consumed, state := n.onEvent(e)
		return consumed, int64(state)
	}
	model.NativeFuncs["gomatcha.io/matcha/pointer Reset"] = func() {
		n.reset()
	}
}

func (r *touchMiddleware) Key() string {
	return "gomatcha.io/matcha/touch"
}

func (r *touchMiddleware) tick() {
	for _, i := range r.tickNodes {
		layoutGuide := r.root.LayoutGuide(int64(i.id))
		e := &event{
			Timestamp: time.Now(),
			Phase:     phaseNone,
			ViewSize:  layoutGuide.Frame.Size(),
		}
		i.onEvent(e)
	}
}

func idSliceToIntSlice(ids []view.Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

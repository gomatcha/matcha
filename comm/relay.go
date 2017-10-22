// Package comm provides primitives for signaling when an object has changed.
// This is done through the Notifier interface and its related children.
package comm

import (
	"strconv"
	"sync"
)

// Relay implements the Notifier interface, and provides methods for triggering notifications
// and republishing notifications from other notifiers. It can be embeded into other structs.
type Relay struct {
	mu    sync.Mutex
	subs  map[Notifier]Id
	funcs map[Id]func()
	maxId Id
}

// Subscribes to notifications from n. When n posts a notification, r will also post a
// notification. For every Subscribe there should be a corresponding Unsubscribe
// or memory leaks may occur.
func (r *Relay) Subscribe(n Notifier) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Multiple subscriptions on the same object are ignored.
	if _, ok := r.subs[n]; ok {
		return
	}

	id := n.Notify(func() {
		r.mu.Lock()
		defer r.mu.Unlock()

		for _, f := range r.funcs {
			f()
		}
	})
	if r.subs == nil {
		r.subs = map[Notifier]Id{}
	}
	r.subs[n] = id
}

// Unsubscribes from n.
func (r *Relay) Unsubscribe(n Notifier) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, ok := r.subs[n]
	if !ok {
		return
	}
	n.Unnotify(id)
	delete(r.subs, n)
}

// Notify implements the Notifier interface.
func (r *Relay) Notify(f func()) Id {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.funcs == nil {
		r.funcs = map[Id]func(){}
	}
	r.maxId += 1
	r.funcs[r.maxId] = f
	return r.maxId
}

// Unnotify implements the Notifier interface.
func (r *Relay) Unnotify(id Id) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.funcs == nil {
		r.funcs = map[Id]func(){}
	}
	if _, ok := r.funcs[id]; !ok {
		panic("comm.Unnotify(): on unknown id: " + strconv.Itoa(int(id)))
	}
	delete(r.funcs, id)
}

// Signal causes all Notifiers on r to be triggered.
func (r *Relay) Signal() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, f := range r.funcs {
		f()
	}
}

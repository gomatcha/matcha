package comm

import (
	"sync"
)

type Relay struct {
	mu    sync.Mutex
	subs  map[Notifier]Id
	funcs map[Id]func()
	maxId Id
}

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

func (r *Relay) Unnotify(id Id) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.funcs == nil {
		r.funcs = map[Id]func(){}
	}
	if _, ok := r.funcs[id]; !ok {
		panic("comm.Unnotify(): on unknown id")
	}
	delete(r.funcs, id)
}

func (r *Relay) Signal() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, f := range r.funcs {
		f()
	}
}

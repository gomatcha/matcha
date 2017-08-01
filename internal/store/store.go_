package store

import (
	"sync"

	"gomatcha.io/matcha/comm"
)

type Storer interface {
	StoreNode() *Node
}

type Node struct {
	mu       sync.Mutex
	mu2      sync.Mutex
	locked   bool
	children map[string]*Node
	parent   *Node

	updated bool

	funcsMu sync.Mutex
	funcs   map[comm.Id]func()
	// notifiers []keyNotifier
	// keyNotifiers map[keyNotifierKey]struct{}
	maxId comm.Id
}

func (s *Node) Set(key string, chl Storer) {
	if s.children == nil {
		s.children = map[string]*Node{}
	}

	// Remove old child.
	s.Delete(key)

	// Add child.
	chlStore := chl.StoreNode()
	if s.locked {
		chlStore.lock()
	}
	s.children[key] = chlStore
	chlStore.parent = s
}

func (s *Node) Delete(key string) {
	if s.children == nil {
		s.children = map[string]*Node{}
	}

	// Flag as updated.
	s.Signal()

	// Remove child.
	if prevChild, ok := s.children[key]; ok {
		if s.locked {
			prevChild.unlock()
		}
		prevChild.parent = nil
	}
	delete(s.children, key)
}

func (s *Node) root() *Node {
	n := s
	for n.parent != nil {
		n = n.parent
	}
	return n
}

func (s *Node) Lock() {
	s.mu.Lock()
	s.root().lock()
}

func (s *Node) lock() {
	s.mu2.Lock()
	for _, i := range s.children {
		i.lock()
	}

	s.locked = true
	s.updated = false
}

func (s *Node) Unlock() {
	s.root().unlock()
	s.mu.Unlock()
}

func (s *Node) unlock() {
	go func(updated bool) {
		s.funcsMu.Lock()
		defer s.funcsMu.Unlock()

		if updated {
			for _, i := range s.funcs {
				i()
			}
		}
	}(s.updated)
	s.updated = false
	s.locked = false

	s.mu2.Unlock()
	for _, i := range s.children {
		i.unlock()
	}
}

func (s *Node) Notify(f func()) comm.Id {
	s.funcsMu.Lock()
	defer s.funcsMu.Unlock()

	if s.funcs == nil {
		s.funcs = map[comm.Id]func(){}
	}

	s.maxId += 1
	s.funcs[s.maxId] = f
	return s.maxId
}

func (s *Node) Unnotify(id comm.Id) {
	s.funcsMu.Lock()
	defer s.funcsMu.Unlock()

	if _, ok := s.funcs[id]; !ok {
		panic("store.Unnotify(): Unknown id")
	}
	delete(s.funcs, id)
}

// func (s *Node) keyNotify(key string, kn *keyNotifier) {
// 	s.funcsMu.Lock()
// 	defer s.funcsMu.Unlock()

// 	if s.keyNotifiers == nil {
// 		s.keyNotifiers = map[keyNotifierKey]struct{}{}
// 	}
// 	s.keyNotifiers[keyNotifierKey{key: key, kn: kn}] = struct{}{}
// }

// func (s *Node) keyUnnotify(key string, kn *keyNotifier) {
// }

// func (s *Node) Notifier(keys ...string) *Notifier {
// 	return keyNotifier{keys: keys, n: s}
// }

func (s *Node) StoreNode() *Node {
	return s
}

func (s *Node) Signal() { // key string
	if s.parent != nil {
		s.parent.updated = true
	}
	s.updated = true
}

// type keyNotifierKey struct {
// 	key string
// 	kn  *keyNotifier
// }

// type keyNotifier struct {
// 	keys []string
// 	n    *Node
// }

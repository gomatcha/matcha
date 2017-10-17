package internal

import (
	"sync"

	"gomatcha.io/matcha/bridge"
)

var tickers = struct {
	ts     map[int]*Ticker
	mu     *sync.Mutex
	maxKey int
}{
	ts:     map[int]*Ticker{},
	mu:     &sync.Mutex{},
	maxKey: 0,
}

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/animate screenUpdate", screenUpdate)
}

func screenUpdate() {
	tickers.mu.Lock()
	ts := []*Ticker{}
	for _, i := range tickers.ts {
		ts = append(ts, i)
	}
	tickers.mu.Unlock()

	for _, i := range ts {
		i.signal()
	}
}

type Ticker struct {
	key int
	mu  sync.Mutex
	f   func()
}

func NewTicker(f func()) *Ticker {
	tickers.mu.Lock()
	defer tickers.mu.Unlock()

	tickers.maxKey += 1
	t := &Ticker{
		key: tickers.maxKey,
		f:   f,
	}
	tickers.ts[t.key] = t
	return t
}

func (t *Ticker) Stop() {
	tickers.mu.Lock()
	defer tickers.mu.Unlock()
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(tickers.ts, t.key)
}

func (t *Ticker) signal() {
	t.mu.Lock()
	f := t.f
	t.mu.Unlock()

	f()
}

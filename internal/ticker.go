package internal

import (
	"sync"
	"time"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/comm"
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
		i.Signal()
	}
}

type Ticker struct {
	key      int
	mu       sync.Mutex
	funcs    map[comm.Id]func()
	maxId    comm.Id
	timer    *time.Timer
	start    time.Time
	duration time.Duration
}

func NewTicker(duration time.Duration) *Ticker {
	tickers.mu.Lock()
	defer tickers.mu.Unlock()

	tickers.maxKey += 1
	t := &Ticker{
		key:      tickers.maxKey,
		funcs:    map[comm.Id]func(){},
		start:    time.Now(),
		duration: duration,
	}
	t.timer = time.AfterFunc(duration, func() {
		t.Stop()
	})
	tickers.ts[t.key] = t
	return t
}

func (t *Ticker) Notify(f func()) comm.Id {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.maxId += 1
	t.funcs[t.maxId] = f
	return t.maxId
}

func (t *Ticker) Unnotify(key comm.Id) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.funcs, key)
}

func (t *Ticker) Value() float64 {
	v := float64(time.Since(t.start)) / float64(t.duration)
	if v < 0 {
		v = 0
	} else if v > 1 {
		v = 1
	}
	return v
}

func (t *Ticker) Stop() {
	tickers.mu.Lock()
	defer tickers.mu.Unlock()
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(tickers.ts, t.key)
}

func (t *Ticker) Signal() {
	t.mu.Lock()
	funcs := t.funcs
	t.mu.Unlock()

	for _, f := range funcs {
		f()
	}
}

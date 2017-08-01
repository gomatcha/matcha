package internal

// func TestScreenUpdate(t *testing.T) {
// 	mu := sync.Mutex{}
// 	count := 0

// 	ticker := NewTicker(time.Second / 20)
// 	done := matcha.NotifyFunc(ticker, func() {
// 		mu.Lock()
// 		defer mu.Unlock()
// 		count += 1
// 	})

// 	screenSignal()
// 	screenSignal()
// 	screenSignal()

// 	mu.Lock()
// 	if count != 3 {
// 		t.Error("Ticker did not trigger")
// 	}
// 	mu.Unlock()

// 	<-time.After(time.Second / 10)

// 	screenSignal()
// 	screenSignal()
// 	screenSignal()

// 	if count != 3 {
// 		t.Error("Ticker not stopped")
// 	}

// 	ticker.Stop()
// 	close(done)
// }

// func TestTickerNotify(t *testing.T) {
// 	ticker := NewTicker(time.Second * 10)
// 	c1 := ticker.Notify()
// 	c2 := ticker.Notify()
// 	c3 := ticker.Notify()

// 	_, ok1 := ticker.chans[c1]
// 	_, ok2 := ticker.chans[c2]
// 	_, ok3 := ticker.chans[c3]

// 	if len(ticker.chans) != 3 || !ok1 || !ok2 || !ok3 {
// 		t.Error("Channel not added")
// 	}

// 	ticker.Unnotify(c1)
// 	ticker.Unnotify(c2)
// 	ticker.Unnotify(c3)

// 	if len(ticker.chans) != 0 {
// 		t.Error("Channel not removed")
// 	}

// 	screenSignal()

// 	ticker.Stop()
// }

// func TestNewTicker(t *testing.T) {
// 	prevCount := len(tickers.ts)

// 	ticker := NewTicker(time.Second * 10)

// 	if len(tickers.ts) != prevCount+1 {
// 		t.Error("Ticker not added")
// 	}
// 	if tickers.maxKey != ticker.key {
// 		t.Error("MaxKey does not match ticker key")
// 	}
// 	if tickers.ts[ticker.key] != ticker {
// 		t.Error("Ticker not found in tickers")
// 	}

// 	ticker.Stop()

// 	if len(tickers.ts) != prevCount {
// 		t.Error("Ticker not removed")
// 	}
// 	if _, ok := tickers.ts[ticker.key]; ok {
// 		t.Error("Ticker not removed from tickers")
// 	}
// }

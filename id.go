package matcha // import "gomatcha.io/matcha"

import (
	"sync"

	_ "gomatcha.io/bridge"
)

var MainLocker sync.Locker

func init() {
	MainLocker = &sync.Mutex{}
}

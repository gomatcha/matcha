package matcha // import "gomatcha.io/matcha"

import (
	"sync"

	_ "gomatcha.io/bridge"
)

type Id int64

var MainLocker sync.Locker

func init() {
	MainLocker = &sync.Mutex{}
}

package matcha

import (
	"sync"

	_ "github.com/gomatcha/matcha/bridge"
)

var MainLocker sync.Locker

func init() {
	MainLocker = &sync.Mutex{}
}

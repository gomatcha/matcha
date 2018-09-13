package device

import (
	"github.com/gomatcha/matcha/bridge"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/internal/device setScreenScale", func(v float64) {
		ScreenScale = v
	})
}

var ScreenScale = 1.0

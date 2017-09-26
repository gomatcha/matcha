package device

import (
	"gomatcha.io/matcha/bridge"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/internal/device setScreenScale", func(v float64) {
		ScreenScale = v
	})
}

var ScreenScale = 1.0

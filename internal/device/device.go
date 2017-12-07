package device

import (
	"math"

	"gomatcha.io/matcha/bridge"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/internal/device setScreenScale", func(v float64) {
		ScreenScale = v
	})
}

var ScreenScale = 1.0

func RoundToScreenScale(v float64) float64 {
	return math.Floor(v*ScreenScale+0.5) / ScreenScale
}

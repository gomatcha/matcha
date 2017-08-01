package animate

import (
	"image/color"

	"gomatcha.io/matcha/comm"
)

// ColorInterpolater represents an object that interpolates between colors given a float64 between 0-1.
type ColorInterpolater interface {
	Interpolate(float64) color.Color
}

// ColorInterpolate wraps n and returns a notifier with the corresponding interpolated colors.
func ColorInterpolate(n comm.Float64Notifier, l ColorInterpolater) comm.ColorNotifier {
	return &colorInterpolater{
		watcher:      n,
		interpolater: l,
	}
}

type colorInterpolater struct {
	watcher      comm.Float64Notifier
	interpolater ColorInterpolater
}

func (w *colorInterpolater) Notify(f func()) comm.Id {
	return w.watcher.Notify(f)
}

func (w *colorInterpolater) Unnotify(id comm.Id) {
	w.watcher.Unnotify(id)
}

func (w *colorInterpolater) Value() color.Color {
	return w.interpolater.Interpolate(w.watcher.Value())
}

// RGBALerp interpolates between colors Start and End.
type RGBALerp struct {
	Start, End color.Color
}

// Interpolate implements the ColorInterpolater interface
func (e RGBALerp) Interpolate(a float64) color.Color {
	r1, g1, b1, a1 := e.Start.RGBA()
	r2, g2, b2, a2 := e.End.RGBA()
	color := color.RGBA64{
		R: uintInterpolate(r1, r2, a),
		G: uintInterpolate(g1, g2, a),
		B: uintInterpolate(b1, b2, a),
		A: uintInterpolate(a1, a2, a),
	}
	return color
}

// Notifier is a convenience method around animate.ColorInterpolate(n, e)
func (e RGBALerp) Notifier(n comm.Float64Notifier) comm.ColorNotifier {
	return ColorInterpolate(n, e)
}

func uintInterpolate(a, b uint32, c float64) uint16 {
	return uint16(a + uint32(float64(b-a)*c))
}

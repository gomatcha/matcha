package animate

import (
	"math"

	"golang.org/x/mobile/exp/sprite/clock"
	"gomatcha.io/matcha/comm"
)

// ColorInterpolater represents an object that interpolates between floats given a float64 between 0-1.
type FloatInterpolater interface {
	Interpolate(float64) float64
}

// FloatInterpolate wraps n and returns a notifier with the corresponding interpolated floats.
func FloatInterpolate(w comm.Float64Notifier, l FloatInterpolater) comm.Float64Notifier {
	return &floatInterpolater{
		watcher:      w,
		interpolater: l,
	}
}

type floatInterpolater struct {
	watcher      comm.Float64Notifier
	interpolater FloatInterpolater
}

func (w *floatInterpolater) Notify(f func()) comm.Id {
	return w.watcher.Notify(f)
}

func (w *floatInterpolater) Unnotify(id comm.Id) {
	w.watcher.Unnotify(id)
}

func (w *floatInterpolater) Value() float64 {
	return w.interpolater.Interpolate(w.watcher.Value())
}

var (
	DefaultEase      FloatInterpolater = CubicBezierEase{0.25, 0.1, 0.25, 1}
	DefaultInEase    FloatInterpolater = CubicBezierEase{0.42, 0, 1, 1}
	DefaultOutEase   FloatInterpolater = CubicBezierEase{0, 0, 0.58, 1}
	DefaultInOutEase FloatInterpolater = CubicBezierEase{0.42, 0, 0.58, 1}
)

// CubicBezierEase interpolates between 1-0 using a Cubic Bézier curve. The parameters are cubic control parameters. The curve starts at (0,0) going toward (x0,y0), and arrives at (1,1) coming from (x1,y1).
type CubicBezierEase struct {
	X0, Y0, X1, Y1 float64
}

// Interpolate implements the Interpolater interface.
func (e CubicBezierEase) Interpolate(a float64) float64 {
	f := clock.CubicBezier(float32(e.X0), float32(e.Y0), float32(e.X1), float32(e.Y1))
	t := f(0, 100000, clock.Time(a*100000))
	return float64(t)
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e).
func (e CubicBezierEase) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

// LinearEase interpolates between 1-0 in a linear fashion.
type LinearEase struct {
}

// Interpolate implements the Interpolater interface.
func (e LinearEase) Interpolate(a float64) float64 {
	return a
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e)
func (e LinearEase) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

// PolyInEase interpolates between Start and End with a polynomial easing.
type PolyInEase struct {
	Exp float64
}

// Interpolate implements the Interpolater interface.
func (e PolyInEase) Interpolate(a float64) float64 {
	return math.Pow(a, e.Exp)
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e)
func (e PolyInEase) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

// PolyOutEase interpolates between Start and End with a reverse polynomial easing.
type PolyOutEase struct {
	Exp float64
}

// Interpolate implements the Interpolater interface.
func (e PolyOutEase) Interpolate(a float64) float64 {
	return 1 - math.Pow(1-a, e.Exp)
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e)
func (e PolyOutEase) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

// PolyInOutEase interpolates between Start and End with a symmetric polynomial easing.
type PolyInOutEase struct {
	ExpIn  float64
	ExpOut float64
}

// Interpolate implements the Interpolater interface.
func (e PolyInOutEase) Interpolate(a float64) float64 {
	if a < 0.5 {
		return math.Pow(a, e.ExpIn)
	} else {
		return 1 - math.Pow(1-a, e.ExpOut)
	}
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e)
func (e PolyInOutEase) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

// FloatLerp interpolates between Start and End linearly.
type FloatLerp struct {
	Start, End float64
}

// Interpolate implements the Interpolater interface.
func (f FloatLerp) Interpolate(a float64) float64 {
	return f.Start + (f.End-f.Start)*a
}

// Notifier is a convenience method around animate.FloatInterpolate(n, e)
func (e FloatLerp) Notifier(a comm.Float64Notifier) comm.Float64Notifier {
	return FloatInterpolate(a, e)
}

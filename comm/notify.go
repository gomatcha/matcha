package comm

import (
	"image/color"
	"time"
)

type Id int64

// Notifier is the interface the wraps the Notify and Unnotify methods.
//
// Notify stores the function f, and calls in the future it when the object updates. It returns
// an Id which can be used to stop notifications. Every call to Notify should have a corresponding Unnotify
// or there could be leaks.
type Notifier interface {
	Notify(f func()) Id
	Unnotify(id Id)
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a color.Color.
type ColorNotifier interface {
	Notifier
	Value() color.Color
}

// ColorNotifier wraps Notifier with an additional Value() method which returns an empty interface.
type InterfaceNotifier interface {
	Notifier
	Value() interface{}
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a bool.
type BoolNotifier interface {
	Notifier
	Value() bool
}

// ColorNotifier wraps Notifier with an additional Value() method which returns an int.
type IntNotifier interface {
	Notifier
	Value() int
}

// ColorNotifier wraps Notifier with an additional Value() method which returns an uint.
type UintNotifier interface {
	Notifier
	Value() uint
}

// ColorNotifier wraps Notifier with an additional Value() method which returns an int64.
type Int64Notifier interface {
	Notifier
	Value() int64
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a uint64.
type Uint64Notifier interface {
	Notifier
	Value() uint64
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a float64.
type Float64Notifier interface {
	Notifier
	Value() float64
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a string.
type StringNotifier interface {
	Notifier
	Value() string
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a []byte.
type BytesNotifier interface {
	Notifier
	Value() []byte
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a time.Duration.
type DurationNotifier interface {
	Notifier
	Value() time.Duration
}

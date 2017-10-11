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

// InterfaceNotifier wraps Notifier with an additional Value() method which returns an empty interface.
type InterfaceNotifier interface {
	Notifier
	Value() interface{}
}

// InterfaceRWNotifier wraps InterfaceNotifier with an additional SetValue(interface{}) method.
type InterfaceRWNotifier interface {
	InterfaceNotifier
	SetValue(interface{})
}

// BoolNotifier wraps Notifier with an additional Value() method which returns a bool.
type BoolNotifier interface {
	Notifier
	Value() bool
}

// BoolRWNotifier wraps BoolNotifier with an additional SetValue(bool) method.
type BoolRWNotifier interface {
	BoolNotifier
	SetValue(bool)
}

// IntNotifier wraps Notifier with an additional Value() method which returns an int.
type IntNotifier interface {
	Notifier
	Value() int
}

// IntRWNotifier wraps IntNotifier with an additional SetValue(int) method.
type IntRWNotifier interface {
	IntNotifier
	SetValue(int)
}

// Float64Notifier wraps Notifier with an additional Value() method which returns a float64.
type Float64Notifier interface {
	Notifier
	Value() float64
}

// Float64RWNotifier wraps Float64Notifier with an additional SetValue(float64) method.
type Float64RWNotifier interface {
	Float64Notifier
	SetValue(float64)
}

// StringNotifier wraps Notifier with an additional Value() method which returns a string.
type StringNotifier interface {
	Notifier
	Value() string
}

// StringRWNotifier wraps StringNotifier with an additional SetValue(string) method.
type StringRWNotifier interface {
	StringNotifier
	SetValue(string)
}

// BytesNotifier wraps Notifier with an additional Value() method which returns a []byte.
type BytesNotifier interface {
	Notifier
	Value() []byte
}

// BytesRWNotifier wraps BytesNotifier with an additional SetValue([]byte) method.
type BytesRWNotifier interface {
	BytesNotifier
	SetValue([]byte)
}

// ColorNotifier wraps Notifier with an additional Value() method which returns a color.Color.
type ColorNotifier interface {
	Notifier
	Value() color.Color
}

// ColorRWNotifier wraps ColorNotifier with an additional SetValue(color.Color) method.
type ColorRWNotifier interface {
	ColorNotifier
	SetValue(color.Color)
}

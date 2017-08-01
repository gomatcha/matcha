package comm

import (
	"image/color"
	"time"
)

type Id int64

type Notifier interface {
	Notify(func()) Id
	Unnotify(Id)
}

type ColorNotifier interface {
	Notifier
	Value() color.Color
}

type InterfaceNotifier interface {
	Notifier
	Value() interface{}
}

type BoolNotifier interface {
	Notifier
	Value() bool
}

type IntNotifier interface {
	Notifier
	Value() int
}

type UintNotifier interface {
	Notifier
	Value() uint
}

type Int64Notifier interface {
	Notifier
	Value() int64
}

type Uint64Notifier interface {
	Notifier
	Value() uint64
}

type Float64Notifier interface {
	Notifier
	Value() float64
}

type StringNotifier interface {
	Notifier
	Value() string
}

type BytesNotifier interface {
	Notifier
	Value() []byte
}

type DurationNotifier interface {
	Notifier
	Value() time.Duration
}

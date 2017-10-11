package comm

import (
	"image/color"
	"sync"
)

// Float64Value implements the Float64RWNotifier interface.
type Float64Value struct {
	value float64
	relay Relay
	mutex sync.Mutex
}

// Notify implements the Float64Notifier interface.
func (v *Float64Value) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the Float64Notifier interface.
func (v *Float64Value) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the Float64Notifier interface.
func (v *Float64Value) Value() float64 {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *Float64Value) SetValue(val float64) {
	v.mutex.Lock()
	if val != v.value {
		v.value = val
		v.mutex.Unlock()
		v.relay.Signal()
	} else {
		v.mutex.Unlock()
	}
}

// IntValue implements the IntRWNotifier interface.
type IntValue struct {
	value int
	relay Relay
	mutex sync.Mutex
}

// Notify implements the Float64Notifier interface.
func (v *IntValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the Float64Notifier interface.
func (v *IntValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the Float64Notifier interface.
func (v *IntValue) Value() int {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *IntValue) SetValue(val int) {
	v.mutex.Lock()
	if val != v.value {
		v.value = val
		v.mutex.Unlock()
		v.relay.Signal()
	} else {
		v.mutex.Unlock()
	}
}

// ColorValue implements the ColorRWNotifier interface.
type ColorValue struct {
	value color.Color
	relay Relay
	mutex sync.Mutex
}

// Notify implements the Float64Notifier interface.
func (v *ColorValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the Float64Notifier interface.
func (v *ColorValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the Float64Notifier interface.
func (v *ColorValue) Value() color.Color {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *ColorValue) SetValue(val color.Color) {
	v.mutex.Lock()
	if val != v.value {
		v.value = val
		v.mutex.Unlock()
		v.relay.Signal()
	} else {
		v.mutex.Unlock()
	}
}

// InterfaceValue implements the InterfaceRWNotifier interface.
type InterfaceValue struct {
	value interface{}
	relay Relay
	mutex sync.Mutex
}

// Notify implements the Float64Notifier interface.
func (v *InterfaceValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the Float64Notifier interface.
func (v *InterfaceValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the Float64Notifier interface.
func (v *InterfaceValue) Value() interface{} {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *InterfaceValue) SetValue(val interface{}) {
	v.mutex.Lock()
	if val != v.value {
		v.value = val
		v.mutex.Unlock()
		v.relay.Signal()
	} else {
		v.mutex.Unlock()
	}
}

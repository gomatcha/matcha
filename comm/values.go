package comm

import (
	"image/color"
	"sync"
)

// InterfaceValue implements the InterfaceRWNotifier interface.
type InterfaceValue struct {
	value interface{}
	relay Relay
	mutex sync.Mutex
}

// Notify implements the InterfaceNotifier interface.
func (v *InterfaceValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the InterfaceNotifier interface.
func (v *InterfaceValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value returns the interface{} stored in v.
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

// BoolValue implements the BoolRWNotifier interface.
type BoolValue struct {
	value bool
	relay Relay
	mutex sync.Mutex
}

// Notify implements the BoolNotifier interface.
func (v *BoolValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the BoolNotifier interface.
func (v *BoolValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the BoolNotifier interface.
func (v *BoolValue) Value() bool {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *BoolValue) SetValue(val bool) {
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

// Notify implements the IntNotifier interface.
func (v *IntValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the IntNotifier interface.
func (v *IntValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the IntNotifier interface.
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

// StringValue implements the StringRWNotifier interface.
type StringValue struct {
	value string
	relay Relay
	mutex sync.Mutex
}

// Notify implements the StringNotifier interface.
func (v *StringValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the StringNotifier interface.
func (v *StringValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value returns the string stored in v.
func (v *StringValue) Value() string {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *StringValue) SetValue(val string) {
	v.mutex.Lock()
	if val != v.value {
		v.value = val
		v.mutex.Unlock()
		v.relay.Signal()
	} else {
		v.mutex.Unlock()
	}
}

// Bytes implements the BytesRWNotifier interface.
type Bytes struct {
	value []byte
	relay Relay
	mutex sync.Mutex
}

// Notify implements the BytesNotifier interface.
func (v *Bytes) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the BytesNotifier interface.
func (v *Bytes) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value implements the BytesNotifier interface.
func (v *Bytes) Value() []byte {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *Bytes) SetValue(val []byte) {
	v.mutex.Lock()
	v.value = val
	v.mutex.Unlock()
	v.relay.Signal()
}

// ColorValue implements the ColorRWNotifier interface.
type ColorValue struct {
	value color.Color
	relay Relay
	mutex sync.Mutex
}

// Notify implements the ColorRWNotifier interface.
func (v *ColorValue) Notify(f func()) Id {
	return v.relay.Notify(f)
}

// Unnotify implements the ColorRWNotifier interface.
func (v *ColorValue) Unnotify(id Id) {
	v.relay.Unnotify(id)
}

// Value returns the color.Color stored in v.
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

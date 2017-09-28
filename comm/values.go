package comm

import "sync"

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
	defer v.mutex.Unlock()
	if val != v.value {
		v.value = val
		v.relay.Signal()
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
	defer v.mutex.Unlock()
	if val != v.value {
		v.value = val
		v.relay.Signal()
	}
}

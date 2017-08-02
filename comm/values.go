package comm

// Float64Value implements the Float64Notifier interface and a setter that updates the
// value and triggers notifications.
type Float64Value struct {
	value float64
	relay Relay
}

// Convenience function that returns a new Float64Value set to val.
func NewFloat64Value(val float64) *Float64Value {
	v := &Float64Value{}
	v.SetValue(val)
	return v
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
	return v.value
}

// SetValue updates v.Value() and notifies any observers.
func (v *Float64Value) SetValue(val float64) {
	if val != v.value {
		v.value = val
		v.relay.Signal()
	}
}

// +build !matcha

// Package bridge implements the interface between Go and Objective-C. See
// https://gomatcha.io/guide/native-bridge/ for more details.
package bridge

// Value wraps an ObjectiveC object.
type Value struct {
	ref int64
}

// Bridge gets the MatchaObjcBridge singleton, and wraps it in a Value.
func Bridge(a string) *Value {
	return nil
}

// Nil returns the Value representing the ObjectiveC nil.
func Nil() *Value {
	return nil
}

// IsNil returns true if v wraps the ObjectievC nil.
func (v *Value) IsNil() bool {
	return false
}

// Bool creates an NSNumber containing v, and wraps it in a Value.
func Bool(v bool) *Value {
	return nil
}

// ToBool returns v's value expressed as a boolean. v must wrap a NSNumber.
func (v *Value) ToBool() bool {
	return false
}

// Int64 creates an NSNumber containing v, and wraps it in a Value.
func Int64(v int64) *Value {
	return nil
}

// ToInt64 returns v's value expressed as an int64. v must wrap a NSNumber.
func (v *Value) ToInt64() int64 {
	return 0
}

// Float64 creates an NSNumber containing v, and wraps it in a Value.
func Float64(v float64) *Value {
	return nil
}

// ToFloat64 returns v's value expressed as a float64. v must wrap a NSNumber.
func (v *Value) ToFloat64() float64 {
	return 0
}

// String creates an NSString containing v, and wraps it in a Value.
func String(v string) *Value {
	return nil
}

// ToString returns v's value as a string. v must wrap a NSString.
func (v *Value) ToString() string {
	return ""
}

// Bytes creates an NSData containing v, and wraps it in a Value.
func Bytes(v []byte) *Value {
	return nil
}

// ToString returns v's value as a byte slice. v must wrap a NSData.
func (v *Value) ToBytes() []byte {
	return nil
}

// Interface creates an MatchaGoValue containing v, and wraps it in a Value.
func Interface(v interface{}) *Value {
	return nil
}

// ToInterface return's v's value as an interface{}, v must wrap a MatchaGoValue.
func (v *Value) ToInterface() interface{} {
	return nil
}

// Array creates an NSArray containing a, and wraps it in a Value.
func Array(a ...*Value) *Value {
	return nil
}

// ToString returns v's elements a slice of Value. v must wrap a NSArray.
func (v *Value) ToArray() []*Value {
	return nil
}

// Call calls a method on v with signature s and arguments args.
//
// Go:
//  rlt := bridge.Bridge().Call("add::", bridge.Int64(1), bridge.Int64(3))
//  fmt.Printf("1+3=%v", rlt.ToInt64())
// Objective-C:
//  @implementation MatchaObjcBridge (Extensions)
//  - (void)add:(NSNumber *)a :(NSNumber *)b {
//  	return @(a.IntegerValue + b.IntegerValue)
//  }
//  @end
func (v *Value) Call(s string, args ...*Value) *Value {
	return nil
}

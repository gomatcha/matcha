// +build matcha

package bridge

// Go support functions for Objective-C. Note that this
// file is copied into and compiled with the generated
// bindings.

/*
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include "go-foreign.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"runtime"
)

//export TestFunc
func TestFunc() {
	a := Bool(true)
	b := Int64(1234)
	c := Float64(1.234)
	d := String("abc")
	e := Bytes([]byte("def123"))

	fmt.Println("blah", a.ToBool(), b.ToInt64(), c.ToFloat64(), d.ToString(), string(e.ToBytes()), "~")

	arr := Array(a, b, c, d, e)
	fmt.Println("done")
	arr2 := arr.ToArray()

	fmt.Println("blah2", arr2[0].ToBool(), arr2[1].ToInt64(), arr2[2].ToFloat64(), arr2[3].ToString(), string(arr2[4].ToBytes()), "~")

	bridge := Bridge("a")
	fmt.Println("bridge", bridge)
}

type Value struct {
	ref int64
}

func newValue(ref C.FgnRef) *Value {
	v := &Value{ref: int64(ref)}
	if ref != 0 {
		runtime.SetFinalizer(v, func(a *Value) {
			C.MatchaForeignUntrack(a._ref())
		})
	}
	return v
}

func (v *Value) _ref() C.FgnRef {
	return C.FgnRef(v.ref)
}

func Bridge(a string) *Value {
	cstr := cString(a)
	return newValue(C.MatchaForeignBridge(cstr))
}

func Nil() *Value {
	return newValue(C.FgnRef(0))
}

func (v *Value) IsNil() bool {
	return v.ref == 0
}

func Bool(v bool) *Value {
	return newValue(C.MatchaForeignBool(C.bool(v)))
}

func (v *Value) ToBool() bool {
	defer runtime.KeepAlive(v)
	return bool(C.MatchaForeignToBool(v._ref()))
}

func Int64(v int64) *Value {
	return newValue(C.MatchaForeignInt64(C.int64_t(v)))
}

func (v *Value) ToInt64() int64 {
	defer runtime.KeepAlive(v)
	return int64(C.MatchaForeignToInt64(v._ref()))
}

func Float64(v float64) *Value {
	return newValue(C.MatchaForeignFloat64(C.double(v)))
}

func (v *Value) ToFloat64() float64 {
	defer runtime.KeepAlive(v)
	return float64(C.MatchaForeignToFloat64(v._ref()))
}

func String(v string) *Value {
	cstr := cString(v)
	return newValue(C.MatchaForeignString(cstr))
}

func (v *Value) ToString() string {
	defer runtime.KeepAlive(v)
	buf := C.MatchaForeignToString(v._ref())
	return goString(buf)
}

func Bytes(v []byte) *Value {
	cbytes := cBytes(v)
	return newValue(C.MatchaForeignBytes(cbytes))
}

func (v *Value) ToBytes() []byte {
	defer runtime.KeepAlive(v)
	buf := C.MatchaForeignToBytes(v._ref())
	return goBytes(buf)
}

func Interface(v interface{}) *Value {
	// Start with a go value.
	// Reflect on it.
	rv := reflect.ValueOf(v)
	// Track it, turning it into a goref.
	ref := matchaGoTrack(rv)
	// Wrap the goref in an foreign object, returning a foreign ref.
	return newValue(C.MatchaForeignGoRef(ref))
}

func (v *Value) ToInterface() interface{} {
	defer runtime.KeepAlive(v)
	// Start with a foreign ref, referring to a foreign value wrapping a go ref.
	// Get the goref.
	ref := C.MatchaForeignToGoRef(v._ref())
	// Get the go object, and unreflect.
	return matchaGoGet(ref).Interface()
}

func Array(a ...*Value) *Value {
	ref := C.MatchaForeignArray(cArray2(a))
	return newValue(ref)
}

func (v *Value) ToArray() []*Value { // TODO(KD): Untested....
	defer runtime.KeepAlive(v)
	buf := C.MatchaForeignToArray(v._ref())
	return goArray2(buf)
}

// Call accepts `nil` in its variadic arguments
func (v *Value) Call(s string, args ...*Value) *Value {
	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(args)

	return newValue(C.MatchaForeignCall(v._ref(), cString(s), cArray2(args)))
}

func cArray(v []reflect.Value) C.CGoBuffer {
	var cstr C.CGoBuffer
	if len(v) == 0 {
		cstr = C.CGoBuffer{}
	} else {
		buf := new(bytes.Buffer)
		for _, i := range v {
			goref := matchaGoTrack(i)
			err := binary.Write(buf, binary.LittleEndian, goref)
			if err != nil {
				fmt.Println("binary.Write failed:", err)
			}
		}
		cstr = C.CGoBuffer{
			ptr: C.CBytes(buf.Bytes()),
			len: C.int64_t(len(buf.Bytes())),
		}
	}
	return cstr
}

func cArray2(v []*Value) C.CGoBuffer {
	var cstr C.CGoBuffer
	if len(v) == 0 {
		cstr = C.CGoBuffer{}
	} else {
		buf := new(bytes.Buffer)
		for _, i := range v {
			foreignRef := i._ref()
			err := binary.Write(buf, binary.LittleEndian, foreignRef)
			if err != nil {
				fmt.Println("binary.Write failed:", err)
			}
		}
		cstr = C.CGoBuffer{
			ptr: C.CBytes(buf.Bytes()),
			len: C.int64_t(len(buf.Bytes())),
		}
	}
	return cstr
}

func cBytes(v []byte) C.CGoBuffer {
	var cstr C.CGoBuffer
	if len(v) == 0 {
		cstr = C.CGoBuffer{}
	} else {
		cstr = C.CGoBuffer{
			ptr: C.CBytes(v),
			len: C.int64_t(len(v)),
		}
	}
	return cstr
}

func cString(v string) C.CGoBuffer {
	var cstr C.CGoBuffer
	if len(v) == 0 {
		cstr = C.CGoBuffer{}
	} else {
		cstr = C.CGoBuffer{
			ptr: C.CBytes([]byte(v)),
			len: C.int64_t(len(v)),
		}
	}
	return cstr
}

func goArray(buf C.CGoBuffer) []reflect.Value {
	defer C.free(buf.ptr)

	gorefs := make([]int64, buf.len/8)
	err := binary.Read(bytes.NewBuffer(C.GoBytes(buf.ptr, C.int(buf.len))), binary.LittleEndian, gorefs)
	if err != nil {
		panic(err)
	}

	rvs := []reflect.Value{}
	for _, i := range gorefs {
		rv := matchaGoGet(C.GoRef(i))
		rvs = append(rvs, rv)
	}
	return rvs
}

func goArray2(buf C.CGoBuffer) []*Value {
	defer C.free(buf.ptr)

	fgnRef := make([]int64, buf.len/8)
	err := binary.Read(bytes.NewBuffer(C.GoBytes(buf.ptr, C.int(buf.len))), binary.LittleEndian, fgnRef)
	if err != nil {
		panic(err)
	}

	rvs := []*Value{}
	for _, i := range fgnRef {
		rv := newValue(C.FgnRef(i))
		rvs = append(rvs, rv)
	}
	return rvs
}

func goString(buf C.CGoBuffer) string {
	defer C.free(buf.ptr)
	str := C.GoBytes(buf.ptr, C.int(buf.len))
	return string(str)
}

func goBytes(buf C.CGoBuffer) []byte {
	defer C.free(buf.ptr)
	return C.GoBytes(buf.ptr, C.int(buf.len))
}

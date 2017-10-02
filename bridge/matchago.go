// +build matcha

package bridge

// Go support functions for Objective-C. Note that this
// file is copied into and compiled with the generated
// bindings.

/*
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include "matchago.h"
*/
import "C"

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"runtime/debug"
	"sync"
)

var goRoot struct {
	types map[string]reflect.Type
	funcs map[string]reflect.Value
}

func init() {
	goRoot.types = map[string]reflect.Type{}
	goRoot.funcs = map[string]reflect.Value{}

	RegisterFunc("gomatcha.io/matcha/bridge Panic", func() {
		panic("test panic")
	})
	RegisterFunc("gomatcha.io/matcha/bridge Panic2", func() {
		var intptr *int
		*intptr = 0
	})
}

func RegisterType(str string, t reflect.Type) {
	goRoot.types[str] = t
}

func RegisterFunc(str string, f interface{}) {
	goRoot.funcs[str] = reflect.ValueOf(f)
}

//export matchaGoForeign
func matchaGoForeign(v C.ObjcRef) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(newValue(v))
	return matchaGoTrack(rv)
}

//export matchaGoToForeign
func matchaGoToForeign(v C.GoRef) C.ObjcRef {
	defer goRecover()
	val := matchaGoGet(v).Interface().(*Value)
	defer runtime.KeepAlive(val)
	return val._ref()
}

//export matchaGoBool
func matchaGoBool(v C.bool) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(bool(v))
	return matchaGoTrack(rv)
}

//export matchaGoToBool
func matchaGoToBool(v C.GoRef) C.bool {
	defer goRecover()
	return C.bool(matchaGoGet(v).Bool())
}

//export matchaGoInt
func matchaGoInt(v C.int) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(int(v))
	return matchaGoTrack(rv)
}

//export matchaGoInt64
func matchaGoInt64(v C.int64_t) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(int64(v))
	return matchaGoTrack(rv)
}

//export matchaGoToInt64
func matchaGoToInt64(v C.GoRef) C.int64_t {
	defer goRecover()
	return C.int64_t(matchaGoGet(v).Int())
}

//export matchaGoUint64
func matchaGoUint64(v C.uint64_t) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(uint64(v))
	return matchaGoTrack(rv)
}

//export matchaGoToUint64
func matchaGoToUint64(v C.GoRef) C.uint64_t {
	defer goRecover()
	return C.uint64_t(matchaGoGet(v).Uint())
}

//export matchaGoFloat64
func matchaGoFloat64(v C.double) C.GoRef {
	defer goRecover()
	rv := reflect.ValueOf(float64(v))
	return matchaGoTrack(rv)
}

//export matchaGoToFloat64
func matchaGoToFloat64(v C.GoRef) C.double {
	defer goRecover()
	return C.double(matchaGoGet(v).Float())
}

//export matchaGoString
func matchaGoString(v C.CGoBuffer) C.GoRef {
	defer goRecover()
	str := goString(v)
	rv := reflect.ValueOf(str)
	return matchaGoTrack(rv)
}

//export matchaGoToString
func matchaGoToString(v C.GoRef) C.CGoBuffer {
	defer goRecover()
	str := matchaGoGet(v).String()
	return C.CGoBuffer{
		ptr: C.CBytes([]byte(str)),
		len: C.int64_t(len(str)),
	}
}

//export matchaGoBytes
func matchaGoBytes(v C.CGoBuffer) C.GoRef {
	defer goRecover()
	defer C.free(v.ptr)
	bytes := C.GoBytes(v.ptr, C.int(v.len))
	rv := reflect.ValueOf(bytes)
	return matchaGoTrack(rv)
}

//export matchaGoToBytes
func matchaGoToBytes(v C.GoRef) C.CGoBuffer {
	defer goRecover()
	bytes := matchaGoGet(v).Bytes()
	return C.CGoBuffer{
		ptr: C.CBytes([]byte(bytes)),
		len: C.int64_t(len(bytes)),
	}
}

//export matchaGoArray
func matchaGoArray() C.GoRef {
	defer goRecover()
	array := []reflect.Value{}
	rv := reflect.ValueOf(array)
	return matchaGoTrack(rv)
}

//export matchaGoArrayLen
func matchaGoArrayLen(v C.GoRef) C.int64_t {
	defer goRecover()
	array := matchaGoGet(v).Interface().([]reflect.Value)
	return C.int64_t(len(array))
}

//export matchaGoArrayAppend
func matchaGoArrayAppend(v, a C.GoRef) C.GoRef {
	defer goRecover()
	array := matchaGoGet(v).Interface().([]reflect.Value)
	elem := matchaGoGet(a)
	newArray := append(array, elem)
	rv := reflect.ValueOf(newArray)
	return matchaGoTrack(rv)
}

//export matchaGoArrayAt
func matchaGoArrayAt(v C.GoRef, idx C.int64_t) C.GoRef {
	defer goRecover()
	array := matchaGoGet(v).Interface().([]reflect.Value)
	return matchaGoTrack(array[idx])
}

//export matchaGoArrayBuffer
func matchaGoArrayBuffer(v C.GoRef) C.CGoBuffer {
	defer goRecover()
	array := matchaGoGet(v).Interface().([]reflect.Value)
	return cArray(array)
}

//export matchaGoMap
func matchaGoMap() C.GoRef {
	defer goRecover()
	m := map[reflect.Value]reflect.Value{}
	rv := reflect.ValueOf(m)
	return matchaGoTrack(rv)
}

//export matchaGoMapKeys
func matchaGoMapKeys(v C.GoRef) C.GoRef {
	defer goRecover()
	keys := matchaGoGet(v).MapKeys()
	return matchaGoTrack(reflect.ValueOf(keys))
}

//export matchaGoMapGet
func matchaGoMapGet(v, key C.GoRef) C.GoRef {
	defer goRecover()
	m := matchaGoGet(v)
	k := matchaGoGet(key)
	return matchaGoTrack(m.MapIndex(k))
}

//export matchaGoMapSet
func matchaGoMapSet(m, key, value C.GoRef) {
	defer goRecover()
	matchaGoGet(m).SetMapIndex(matchaGoGet(key), matchaGoGet(value))
}

//export matchaGoType
func matchaGoType(v C.CGoBuffer) C.GoRef {
	defer goRecover()
	str := goString(v)
	t := goRoot.types[str]
	rv := reflect.New(t)
	return matchaGoTrack(rv)
}

//export matchaGoFunc
func matchaGoFunc(v C.CGoBuffer) C.GoRef {
	defer goRecover()
	str := goString(v)
	f, ok := goRoot.funcs[str]
	if !ok {
		fmt.Println("No such function:", str)
	}
	return matchaGoTrack(f)
}

//export matchaGoIsNil
func matchaGoIsNil(v C.GoRef) C.bool {
	defer goRecover()
	return C.bool(matchaGoGet(v).IsNil())
}

//export matchaGoEqual
func matchaGoEqual(a C.GoRef, b C.GoRef) C.bool {
	defer goRecover()
	return C.bool(matchaGoGet(a).Interface() == matchaGoGet(b).Interface())
}

//export matchaGoElem
func matchaGoElem(v C.GoRef) C.GoRef {
	defer goRecover()
	rv := matchaGoGet(v)
	return matchaGoTrack(rv.Elem())
}

//export matchaGoCall
func matchaGoCall(v C.GoRef, name C.CGoBuffer, args C.GoRef) C.GoRef {
	defer goRecover()
	str := goString(name)
	rv := matchaGoGet(v)

	var function reflect.Value
	if str == "" {
		function = rv
	} else {
		function = rv.MethodByName(str)
	}
	argsRv := matchaGoGet(args).Interface().([]reflect.Value)

	rlt := function.Call(argsRv)
	return matchaGoTrack(reflect.ValueOf(rlt))
}

//export matchaGoCall2
func matchaGoCall2(v C.GoRef, name C.CGoBuffer, args C.CGoBuffer) C.CGoBuffer {
	defer goRecover()
	str := goString(name)
	rv := matchaGoGet(v)

	var function reflect.Value
	if str == "" {
		function = rv
	} else {
		function = rv.MethodByName(str)
	}
	argsRv := goArray(args)
	rlt := function.Call(argsRv)
	return cArray(rlt)
}

//export matchaGoField
func matchaGoField(v C.GoRef, name C.CGoBuffer) C.GoRef {
	defer goRecover()
	rv := matchaGoGet(v)
	str := goString(name)

	// Always underlying value.
	kind := rv.Kind()
	for kind == reflect.Ptr || kind == reflect.Interface {
		rv = rv.Elem()
		kind = rv.Kind()
	}

	field := rv.FieldByName(str)
	return matchaGoTrack(field)
}

//export matchaGoFieldSet
func matchaGoFieldSet(v C.GoRef, name C.CGoBuffer, elem C.GoRef) {
	defer goRecover()
	rv := matchaGoGet(v)
	str := goString(name)

	// Always underlying value.
	kind := rv.Kind()
	for kind == reflect.Ptr || kind == reflect.Interface {
		rv = rv.Elem()
		kind = rv.Kind()
	}

	rv.FieldByName(str).Set(matchaGoGet(elem))
}

var tracker struct {
	sync.Mutex
	minRef int64
	refs   map[int64]reflect.Value
}

func init() {
	tracker.refs = map[int64]reflect.Value{}
}

func matchaGoTrack(v reflect.Value) C.GoRef {
	tracker.Lock()
	defer tracker.Unlock()
	if len(tracker.refs)%100 == 0 {
		fmt.Println("track count:", len(tracker.refs))
	}

	tracker.minRef -= 1
	tracker.refs[tracker.minRef] = v
	return C.GoRef(tracker.minRef)
}

func matchaGoGet(ref C.GoRef) reflect.Value {
	tracker.Lock()
	defer tracker.Unlock()

	v, ok := tracker.refs[int64(ref)]
	if !ok {
		panic("Get error. No corresponding object for key.")
	}
	return v
}

//export matchaGoUntrack
func matchaGoUntrack(ref C.GoRef) {
	defer goRecover()
	tracker.Lock()
	defer tracker.Unlock()

	_, ok := tracker.refs[int64(ref)]
	if !ok {
		panic("Untrack error. No corresponding object for key.")
	}
	delete(tracker.refs, int64(ref))
}

// For better crash logs on Android
func goRecover() {
	if r := recover(); r != nil {
		log.Printf("%s %s", r, debug.Stack())
		C.MatchaForeignPanic()
	}
}

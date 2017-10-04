// +build matcha,android

package bridge

/*
#cgo CFLAGS:
#cgo LDFLAGS: -landroid -llog

#include "go-foreign.h"
#include "go-go.h"
#include "java-foreign.h"
#include "java-go.h"
*/
import "C"

// +build matcha,darwin

package bridge

/*
#cgo CFLAGS: -x objective-c -fobjc-arc -Werror
#cgo LDFLAGS: -framework Foundation

#include "objc-go.h"
#include "objc-foreign.h"
*/
import "C"

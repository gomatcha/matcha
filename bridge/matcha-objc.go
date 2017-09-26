// +build matcha,darwin

package bridge

/*
#cgo CFLAGS: -x objective-c -fobjc-arc -Werror
#cgo LDFLAGS: -framework Foundation

#include "matchago-objc.h"
#include "matchaforeign-objc.h"
*/
import "C"

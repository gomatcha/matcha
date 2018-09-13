/*
Package application provides access to application resources. Image assets
must be in the app's .xcassets file (iOS) or res/drawable folder (Android).
Disable "Compress PNG Files" and "Remove Text Metadata from PNG Files" in Xcode
if loading image resources is not working. Android does not allow uppercase
image names or folders and this restriction carries over to Matcha as well.

    // Display an image.
    img, err := application.LoadImage("example")
    if err != nil {
        imageview.Image = img
    }

    // or
    imageview.Image = application.MustLoadImage("example")
*/
package application

import (
	"errors"
	"runtime"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/layout"
)

// // AssetsDir returns the path to the app's assets directory. `NSBundle.mainBundle.resourcePath`
// func AssetsDir() (string, error) {
// 	return bridge.Bridge("").Call("assetsDir").ToString(), nil
// }

func OpenURL(url string) error {
	success := true
	if runtime.GOOS == "android" {
		success = bridge.Bridge("").Call("openURL", bridge.String(url)).ToBool()
	} else {
		success = bridge.Bridge("").Call("openURL:", bridge.String(url)).ToBool()
	}
	if !success {
		return errors.New("Unable to open URL")
	}
	return nil
}

func orientation(v int) layout.Edge {
	switch v {
	case 0:
		return layout.EdgeTop
	case 1:
		return layout.EdgeBottom
	case 2:
		return layout.EdgeRight
	case 3:
		return layout.EdgeLeft
	}
	return layout.EdgeTop
}

// layout.EdgeTop is portrait. EdgeRight and EdgeLeft are landscape and
// EdgeBottom is upside down.
func Orientation() layout.Edge {
	var o int64
	if runtime.GOOS == "android" {
		o = bridge.Bridge("").Call("orientation").ToInt64()
	} else {
		o = bridge.Bridge("").Call("orientation").ToInt64()
	}
	return orientation(int(o))
}

// In order to receive orientation notification on android you may need to set
// `android:configChanges="orientation|keyboardHidden|screenSize"` in your
// activity's manifest (https://stackoverflow.com/a/6005624).
var OrientationNotifier comm.IntNotifier
var orientationNotifier comm.IntValue

func init() {
	OrientationNotifier = &orientationNotifier
	bridge.RegisterFunc("github.com/gomatcha/matcha/application SetOrientation", func(v int) {
		orientationNotifier.SetValue(int(orientation(v)))
	})
}

var ShakeNotifier comm.Notifier
var shakeNotifier comm.Relay

func init() {
	ShakeNotifier = &shakeNotifier
	bridge.RegisterFunc("github.com/gomatcha/matcha/application OnShake", func() {
		shakeNotifier.Signal()
	})
}

// Package app provides access to application properties.
package app

import (
	"gomatcha.io/bridge"
)

// AssetsDir returns the path to the app's assets directory. `NSBundle.mainBundle.resourcePath`
func AssetsDir() (string, error) {
	return bridge.Bridge().Call("assetsDir").ToString(), nil
}

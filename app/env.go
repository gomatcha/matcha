package app

import (
	"gomatcha.io/bridge"
)

func AssetsDir() (string, error) {
	return bridge.Bridge().Call("assetsDir").ToString(), nil
}

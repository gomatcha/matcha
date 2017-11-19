package internal

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

var BackRelay *comm.Relay = nil

func IosBackItem() *ios.StackBarItem {
	if BackRelay == nil {
		return nil
	}

	iosItem := ios.NewStackBarItem()
	iosItem.Title = "Examples"
	iosItem.OnPress = func() {
		if BackRelay != nil {
			BackRelay.Signal()
		}
	}
	return iosItem
}

func AndroidBackItem() *android.StackBarItem {
	if BackRelay == nil {
		return nil
	}

	androidItem := android.NewStackBarItem()
	androidItem.Title = "Examples"
	androidItem.OnPress = func() {
		if BackRelay != nil {
			BackRelay.Signal()
		}
	}
	return androidItem
}

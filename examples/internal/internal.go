package internal

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/view/android"
	"gomatcha.io/matcha/view/ios"
)

var ExamplesRelay *comm.Relay = nil

func IosExamplesItem() *ios.StackBarItem {
	if ExamplesRelay == nil {
		return nil
	}

	iosItem := ios.NewStackBarItem()
	iosItem.Title = "Examples"
	iosItem.OnPress = func() {
		if ExamplesRelay != nil {
			ExamplesRelay.Signal()
		}
	}
	return iosItem
}

func AndroidExamplesItem() *android.StackBarItem {
	if ExamplesRelay == nil {
		return nil
	}

	androidItem := android.NewStackBarItem()
	androidItem.Title = "Examples"
	androidItem.OnPress = func() {
		if ExamplesRelay != nil {
			ExamplesRelay.Signal()
		}
	}
	return androidItem
}

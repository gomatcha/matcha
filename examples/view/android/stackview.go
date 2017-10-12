package android

import (
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/android NewStackView", func() view.View {
		return NewStackView()
	})
}

func NewStackView() view.View {
	return ios.NewStackView()
}

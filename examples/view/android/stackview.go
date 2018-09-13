package android

import (
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/examples/view/ios"
	"github.com/gomatcha/matcha/view"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view/android NewStackView", func() view.View {
		return NewStackView()
	})
}

func NewStackView() view.View {
	return ios.NewStackView()
}

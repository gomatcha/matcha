package view

import (
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
)

type BasicView struct {
	Embed
	Children []View
	Layouter layout.Layouter
	Painter  paint.Painter
	Options  []Option
}

// NewBasicView returns a new view.
func NewBasicView() *BasicView {
	return &BasicView{}
}

// Build implements view.View.
func (v *BasicView) Build(ctx Context) Model {
	return Model{
		Children: v.Children,
		Layouter: v.Layouter,
		Painter:  v.Painter,
		Options:  v.Options,
	}
}

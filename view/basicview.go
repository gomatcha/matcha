package view

import (
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
)

type BasicView struct {
	Embed
	Painter  paint.Painter
	Layouter layout.Layouter
	Children []View
}

// NewBasicView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewBasicView() *BasicView {
	return &BasicView{}
}

// Build implements view.View.
func (v *BasicView) Build(ctx *Context) Model {
	return Model{
		Children: v.Children,
		Painter:  v.Painter,
		Layouter: v.Layouter,
	}
}

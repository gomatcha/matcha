// Package basicview implements an empty view.
package basicview

import (
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type View struct {
	view.Embed
	Painter  paint.Painter
	Layouter layout.Layouter
	Children []view.View
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	return view.Model{
		Children: v.Children,
		Painter:  v.Painter,
		Layouter: v.Layouter,
	}
}

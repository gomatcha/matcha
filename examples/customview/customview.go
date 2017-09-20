package customview

import (
	_ "gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type View struct {
	view.Embed
	PaintStyle paint.Painter
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{}
}

// Build implements view.View.
func (v *View) Build(ctx view.Context) view.Model {
	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		NativeViewName: "github.com/overcyn/customview",
		Painter:        painter,
	}
}

package customview

import (
	_ "gomatcha.io/bridge"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type CustomView struct {
	view.Embed
	PaintStyle paint.Painter
}

// New returns an initialized CustomView instance.
func NewCustomView() *View {
	return &CustomView{}
}

// Build implements view.View.
func (v *CustomView) Build(ctx view.Context) view.Model {
	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		NativeViewName: "github.com/overcyn/customview",
		Painter:        painter,
	}
}

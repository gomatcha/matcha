package passthrough

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
)

type Layouter struct {
}

// Layout implements the view.Layouter interface.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := ctx.LayoutChild(0, ctx.MinSize, ctx.MinSize)
	return g, []layout.Guide{g}
}

// Notify implements the view.Layouter interface.
func (l *Layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *Layouter) Unnotify(id comm.Id) {
	// no-op
}

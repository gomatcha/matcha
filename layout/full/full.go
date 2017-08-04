/*
Package full implements a layout system where the view and all direct children are positioned to the maximum size. This is the default layouter. It does not support animations or flexible sizing. For more complex layouts, see the constraint package.

 views := []view.View{}

 childView := NewChildView(...)
 views = append(views, childView)

 return view.Model{
 	Views: views,
 	Layouter:&full.Layouter{},
 }
*/
package full

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
)

type Layouter struct {
}

// Layout implements the view.Layouter interface.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{Frame: layout.Rect{Max: ctx.MinSize}}
	gs := []layout.Guide{}
	for i := range ctx.ChildIds {
		gs = append(gs, ctx.LayoutChildIdx(i, ctx.MinSize, ctx.MinSize))
	}
	return g, gs
}

// Notify implements the view.Layouter interface.
func (l *Layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *Layouter) Unnotify(id comm.Id) {
	// no-op
}

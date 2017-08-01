/*
Package absolute implements a fixed layout system similar to HTML absolute positioning. It does not support animations or flexible sizing. For more complex layouts, see the constraint package.

 l := &absolute.Layouter{
	Guide: layout.Guide{Frame: layout.Rt(0, 0, 100, 100)}
 }

 childView := NewChildView(...)
 l.Add(childView, layout.Guide{Frame: layout.Rt(10, 10, 90, 90)})

 return view.Model{
 	Views: l.Views(),
 	Layouter:l,
 }
*/
package absolute

import (
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/view"
)

type Layouter struct {
	// Layout guide for the view.
	Guide       layout.Guide
	childGuides map[matcha.Id]layout.Guide
	views       []view.View
}

// Add adds v to the layouter and positions it with g.
func (l *Layouter) Add(v view.View, g layout.Guide) {
	if l.childGuides == nil {
		l.childGuides = map[matcha.Id]layout.Guide{}
	}
	l.childGuides[v.Id()] = g
	l.views = append(l.views, v)
}

// Views returns all views that have been added to l.
func (l *Layouter) Views() []view.View {
	return l.views
}

// Layout implements the view.Layouter interface.
func (l *Layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	return l.Guide, l.childGuides
}

// Notify implements the view.Layouter interface.
func (l *Layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *Layouter) Unnotify(id comm.Id) {
	// no-op
}

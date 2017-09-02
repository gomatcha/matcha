package view

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
)

type absoluteLayouter struct {
	// Layout guide for the view.
	Guide       layout.Guide
	childGuides []layout.Guide
	views       []View
}

// Add adds v to the layouter and positions it with g.
func (l *absoluteLayouter) Add(v View, g layout.Guide) {
	l.childGuides = append(l.childGuides, g)
	l.views = append(l.views, v)
}

// Views returns all views that have been added to l.
func (l *absoluteLayouter) Views() []View {
	return l.views
}

// Layout implements the view.Layouter interface.
func (l *absoluteLayouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	// TODO(KD): Need to call layoutChild.
	for i := 0; i < len(l.childGuides); i++ {
		g := l.childGuides[i]
		p := layout.Pt(g.Width(), g.Height())
		ctx.LayoutChild(i, p, p)
	}
	return l.Guide, l.childGuides
}

// Notify implements the view.Layouter interface.
func (l *absoluteLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

// Unnotify implements the view.Layouter interface.
func (l *absoluteLayouter) Unnotify(id comm.Id) {
	// no-op
}

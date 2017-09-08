package android

import "gomatcha.io/matcha/view"

type SwipeView struct {
	view.Embed
}

// NewSwipeView returns either the previous View in ctx with matching key, or a new View if none exists.
// ViewPager and PagerTabStrip.
func NewSwipeView() *SwipeView {
	return &SwipeView{}
}

// Build implements view.View.
func (v *SwipeView) Build(ctx *view.Context) view.Model {
	return view.Model{
		NativeViewName:  "gomatcha.io/matcha/view/android SwipeView",
		NativeViewState: nil,
	}
}

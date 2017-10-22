package view

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
)

// TextView displays a multiline text region within it bounds.
type TextView struct {
	Embed
	PaintStyle *paint.Style
	String     string
	Text       *text.Text
	Style      *text.Style
	StyledText *text.StyledText
	MaxLines   int
	layouter   textViewLayouter
}

// NewTextView returns a new view.
func NewTextView() *TextView {
	return &TextView{
		Style: &text.Style{},
	}
}

// Build implements the view.View interface.
func (v *TextView) Build(ctx Context) Model {
	st := v.StyledText
	if st == nil {
		t := v.Text
		if t == nil {
			t = text.New(v.String)
		}
		st = text.NewStyledText(t.String(), v.Style)
	}

	v.layouter.styledText = st
	v.layouter.maxLines = v.MaxLines

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Painter:         painter,
		Layouter:        &v.layouter,
		NativeViewName:  "gomatcha.io/matcha/view/textview",
		NativeViewState: internal.MarshalProtobuf(st.MarshalProtobuf()),
	}
}

type textViewLayouter struct {
	styledText *text.StyledText
	maxLines   int

	cachedStyledText *text.StyledText
	maxSize          layout.Point
	minSize          layout.Point
	guide            layout.Guide
}

func (l *textViewLayouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	if l.maxSize == ctx.MaxSize() && l.minSize == ctx.MinSize() && l.styledText.Equal(l.cachedStyledText) {
		return l.guide, nil
	}

	size := l.styledText.Size(layout.Pt(0, 0), ctx.MaxSize(), l.maxLines)
	g := layout.Guide{Frame: layout.Rt(0, 0, size.X, size.Y)}

	// Cache values for later
	l.cachedStyledText = l.styledText.Copy()
	l.maxSize = ctx.MaxSize()
	l.minSize = ctx.MinSize()
	l.guide = g

	return g, nil
}

func (l *textViewLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *textViewLayouter) Unnotify(id comm.Id) {
	// no-op
}

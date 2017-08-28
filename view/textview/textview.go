// Package textview implements a text area.
package textview

import (
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

// View displays a multiline text region within it bounds.
type View struct {
	view.Embed
	PaintStyle *paint.Style
	String     string
	Text       *text.Text
	Style      *text.Style
	StyledText *text.StyledText // TODO(KD): subscribe to StyledText and Text
	MaxLines   int
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{
		Style: &text.Style{},
	}
}

// Build implements the view.View interface.
func (v *View) Build(ctx *view.Context) view.Model {
	st := v.StyledText
	if st == nil {
		t := v.Text
		if t == nil {
			t = text.New(v.String)
		}
		st = text.NewStyledText(t.String(), v.Style)
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:         painter,
		Layouter:        &layouter{styledText: st, maxLines: v.MaxLines},
		NativeViewName:  "gomatcha.io/matcha/view/textview",
		NativeViewState: st.MarshalProtobuf(),
	}
}

type layouter struct {
	styledText *text.StyledText
	maxLines   int
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	size := l.styledText.Size(layout.Pt(0, 0), ctx.MaxSize, l.maxLines)
	g := layout.Guide{Frame: layout.Rt(0, 0, size.X, size.Y)}
	return g, nil
}

func (l *layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *layouter) Unnotify(id comm.Id) {
	// no-op
}

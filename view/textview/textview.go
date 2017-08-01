package textview

import (
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

type layouter struct {
	styledText *internal.StyledText
	maxLines   int
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
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

type View struct {
	view.Embed
	PaintStyle *paint.Style
	String     string
	Text       *text.Text
	Style      *text.Style
	MaxLines   int
}

func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed: ctx.NewEmbed(key),
		Style: &text.Style{},
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	t := v.Text
	if t == nil {
		t = text.New(v.String)
	}
	st := internal.NewStyledText(t)
	st.Set(v.Style, 0, 0)

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

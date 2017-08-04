// Package button implements a native button view.
package button

import (
	"image/color"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	pbbutton "gomatcha.io/matcha/pb/view/button"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

// View implements a native button view.
type View struct {
	view.Embed
	Text       string
	Color      color.Color
	OnPress    func()
	Enabled    bool
	PaintStyle *paint.Style
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed:   ctx.NewEmbed(key),
		Enabled: true,
		Color:   color.RGBA{14, 122, 254, 255},
	}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	style := &text.Style{}
	style.SetAlignment(text.AlignmentCenter)
	style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   20,
	})
	style.SetTextColor(v.Color)
	t := text.New(v.Text)
	st := internal.NewStyledText(t)
	st.Set(style, 0, 0)

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:        painter,
		Layouter:       &layouter{styledText: st},
		NativeViewName: "gomatcha.io/matcha/view/button",
		NativeViewState: &pbbutton.View{
			StyledText: st.MarshalProtobuf(),
			Enabled:    v.Enabled,
			Color:      pb.ColorEncode(v.Color),
		},
		NativeFuncs: map[string]interface{}{
			"OnPress": func() {
				if v.OnPress != nil {
					v.OnPress()
				}
			},
		},
	}
}

type layouter struct {
	styledText *internal.StyledText
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	const padding = 10.0
	size := l.styledText.Size(layout.Pt(0, 0), ctx.MaxSize, 1)
	g := layout.Guide{Frame: layout.Rt(0, 0, size.X+padding*2, size.Y+padding*2)}
	return g, nil
}

func (l *layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *layouter) Unnotify(id comm.Id) {
	// no-op
}

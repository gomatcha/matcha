package textinput

import (
	"fmt"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/textinput"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

// View mutates the Text and StyledText fields in place.
type View struct {
	view.Embed
	PaintStyle         *paint.Style
	Text               *text.Text
	text               *text.Text
	Style              *text.Style
	PlaceholderText    *text.Text
	PlaceholderStyle   *text.Style
	SecureTextEntry    bool
	KeyboardType       keyboard.Type
	KeyboardAppearance keyboard.Appearance
	KeyboardReturnType keyboard.ReturnType
	Responder          *keyboard.Responder
	prevResponder      *keyboard.Responder
	responder          *keyboard.Responder
	Multiline          bool
	OnTextChange       func(*text.Text)
	OnSubmit           func()
	OnFocus            func(*keyboard.Responder)
}

func New(ctx *view.Context, key string) *View {
	if v, ok := ctx.Prev(key).(*View); ok {
		return v
	}
	return &View{
		Embed:     ctx.NewEmbed(key),
		text:      text.New(""),
		responder: &keyboard.Responder{},
	}
}

func (v *View) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.prevResponder)
	}
}

func (v *View) Build(ctx *view.Context) view.Model {
	style := v.Style
	if style == nil {
		style = &text.Style{}
	}

	t := v.Text
	if t == nil {
		t = v.text
	}
	st := internal.NewStyledText(t)
	st.Set(style, 0, 0)

	placeholder := v.PlaceholderText
	if placeholder == nil {
		placeholder = text.New("")
	}
	placeholderStyledText := internal.NewStyledText(placeholder)
	placeholderStyledText.Set(v.PlaceholderStyle, 0, 0)

	if v.Responder != v.prevResponder {
		if v.prevResponder != nil {
			v.Unsubscribe(v.prevResponder)
		}

		v.prevResponder = v.Responder
		if v.Responder != nil {
			v.Subscribe(v.Responder)
		}
	}

	responder := v.Responder
	if responder == nil {
		responder = v.responder
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Layouter:       &layouter{styledText: st, multiline: v.Multiline},
		Painter:        painter,
		NativeViewName: "gomatcha.io/matcha/view/textinput",
		NativeViewState: &textinput.View{
			StyledText:         st.MarshalProtobuf(),
			PlaceholderText:    placeholderStyledText.MarshalProtobuf(),
			KeyboardType:       v.KeyboardType.MarshalProtobuf(),
			KeyboardAppearance: v.KeyboardAppearance.MarshalProtobuf(),
			KeyboardReturnType: v.KeyboardReturnType.MarshalProtobuf(),
			Focused:            responder.Visible(),
			Multiline:          v.Multiline,
			SecureTextEntry:    v.SecureTextEntry,
		},
		NativeFuncs: map[string]interface{}{
			"OnTextChange": func(data []byte) {
				pbevent := &textinput.Event{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				_text := v.Text
				if _text == nil {
					_text = v.text
				}

				_text.UnmarshalProtobuf(pbevent.StyledText.Text)
				if v.OnTextChange != nil {
					v.OnTextChange(_text)
				}
			},
			"OnSubmit": func() {
				if v.OnSubmit != nil {
					v.OnSubmit()
				}
			},
			"OnFocus": func(data []byte) {
				pbevent := &textinput.FocusEvent{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				responder := v.Responder
				if responder == nil {
					responder = v.responder
				}

				if pbevent.Focused {
					responder.Show()
				} else {
					responder.Dismiss()
				}
				if v.OnFocus != nil {
					v.OnFocus(responder)
				}
			},
		},
	}
}

type layouter struct {
	styledText *internal.StyledText
	multiline  bool
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, map[matcha.Id]layout.Guide) {
	if !l.multiline {
		size := l.styledText.Size(layout.Pt(0, 0), ctx.MaxSize, 1)
		size.Y += 15
		if size.Y < 30 {
			size.Y = 30
		}
		g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize.X, size.Y)}
		return g, nil
	} else {
		g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize.X, ctx.MinSize.Y)}
		return g, nil
	}
}

func (l *layouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *layouter) Unnotify(id comm.Id) {
	// no-op
}

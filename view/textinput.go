package view

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pbview "gomatcha.io/matcha/pb/view"
	"gomatcha.io/matcha/text"
)

// TextInput represents a text input view. TextInput mutates the Text and
// StyledText fields in place.
type TextInput struct {
	Embed
	PaintStyle         *paint.Style
	Text               *text.Text
	text               *text.Text
	Style              *text.Style
	Placeholder        string
	PlaceholderStyle   *text.Style
	Password           bool
	KeyboardType       keyboard.Type
	KeyboardReturnType keyboard.ReturnType
	Responder          *keyboard.Responder
	prevResponder      *keyboard.Responder
	responder          *keyboard.Responder
	MaxLines           int // This is used only for sizing.
	OnTextChange       func(*text.Text)
	OnSubmit           func()
	OnFocus            func(*keyboard.Responder)
}

// NewTextInput returns either the previous View in ctx with matching key, or a new View if none exists.
func NewTextInput() *TextInput {
	return &TextInput{
		MaxLines:  1,
		text:      text.New(""),
		responder: &keyboard.Responder{},
	}
}

// Lifecyle implements the view.View interface.
func (v *TextInput) Lifecycle(from, to Stage) {
	if ExitsStage(from, to, StageMounted) {
		v.Unsubscribe(v.prevResponder)
	}
}

// Build implements the view.View interface.
func (v *TextInput) Build(ctx Context) Model {
	style := v.Style
	if style == nil {
		style = &text.Style{}
		if runtime.GOOS == "android" {
			style.SetFont(text.DefaultFont(18))
		} else if runtime.GOOS == "darwin" {
			style.SetFont(text.DefaultFont(18))
		}
	}

	t := v.Text
	if t == nil {
		t = v.text
	}
	st := text.NewStyledText(t.String(), style)

	placeholderStyle := v.PlaceholderStyle
	if placeholderStyle == nil {
		placeholderStyle = &text.Style{}
		if runtime.GOOS == "android" {
			placeholderStyle.SetFont(text.DefaultFont(18))
			placeholderStyle.SetTextColor(colornames.Gray)
		} else if runtime.GOOS == "darwin" {
			placeholderStyle.SetFont(text.DefaultFont(18))
			placeholderStyle.SetTextColor(colornames.Lightgray)
		}
	}
	placeholderStyledText := text.NewStyledText(v.Placeholder, placeholderStyle)

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
	return Model{
		Layouter:       &textInputLayouter{style: style, styledText: st, maxLines: v.MaxLines},
		Painter:        painter,
		NativeViewName: "gomatcha.io/matcha/view/textinput",
		NativeViewState: &pbview.TextInput{
			Font:               style.Font().MarshalProtobuf(),
			StyledText:         st.MarshalProtobuf(),
			PlaceholderText:    placeholderStyledText.MarshalProtobuf(),
			KeyboardType:       v.KeyboardType.MarshalProtobuf(),
			KeyboardReturnType: v.KeyboardReturnType.MarshalProtobuf(),
			Focused:            responder.Visible(),
			MaxLines:           int64(v.MaxLines),
			SecureTextEntry:    v.Password,
		},
		NativeFuncs: map[string]interface{}{
			"OnTextChange": func(data []byte) {
				pbevent := &pbview.TextInputEvent{}
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
				pbevent := &pbview.TextInputFocusEvent{}
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

type textInputLayouter struct {
	style      *text.Style
	styledText *text.StyledText
	maxLines   int
}

func (l *textInputLayouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	if l.maxLines == 1 {
		st := text.NewStyledText("AAA", l.style)
		size := st.Size(layout.Pt(0, 0), ctx.MaxSize(), 1)
		g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize().X, size.Y)}
		return g, nil
	} else {
		g := layout.Guide{Frame: layout.Rt(0, 0, ctx.MinSize().X, ctx.MinSize().Y)}
		return g, nil
	}
}

func (l *textInputLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *textInputLayouter) Unnotify(id comm.Id) {
	// no-op
}

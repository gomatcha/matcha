// Package textinput implements a text input field.
package textinput

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb/view/textinput"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

// View represents a text input view. View mutates the Text and
// StyledText fields in place.
type View struct {
	view.Embed
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

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{
		text:      text.New(""),
		responder: &keyboard.Responder{},
	}
}

// Lifecyle implements the view.View interface.
func (v *View) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.prevResponder)
	}
}

// Build implements the view.View interface.
func (v *View) Build(ctx *view.Context) view.Model {
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
	return view.Model{
		Layouter:       &layouter{styledText: st, maxLines: v.MaxLines},
		Painter:        painter,
		NativeViewName: "gomatcha.io/matcha/view/textinput",
		NativeViewState: &textinput.View{
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
	styledText *text.StyledText
	maxLines   int
}

func (l *layouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	if l.maxLines == 1 {
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

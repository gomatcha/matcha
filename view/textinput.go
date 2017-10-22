package view

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pbview "gomatcha.io/matcha/proto/view"
	"gomatcha.io/matcha/text"
)

// TextInput represents a text input view. TextInput mutates the Text and
// StyledText fields in place.
type TextInput struct {
	Embed
	PaintStyle       *paint.Style
	RWText           *text.Text
	Style            *text.Style
	Placeholder      string
	PlaceholderStyle *text.Style
	Password         bool
	KeyboardType     keyboard.Type
	Responder        *keyboard.Responder
	MaxLines         int // This is used only for sizing.
	OnChange         func(*text.Text)
	OnSubmit         func(*text.Text)
	OnFocus          func(*keyboard.Responder)
}

// NewTextInput returns a new view.
func NewTextInput() *TextInput {
	return &TextInput{
		MaxLines: 1,
	}
}

// Lifecyle implements the view.View interface.
func (v *TextInput) Lifecycle(from, to Stage) {
	if EntersStage(from, to, StageMounted) {
		if v.RWText == nil {
			v.RWText = text.New("")
		}
		v.Subscribe(v.RWText)

		if v.Responder == nil {
			v.Responder = &keyboard.Responder{}
		}
		v.Subscribe(v.Responder)
	} else if ExitsStage(from, to, StageMounted) {
		v.Unsubscribe(v.RWText)
		v.Unsubscribe(v.Responder)
	}
}

func (v *TextInput) Update(v2 View) {
	v.Unsubscribe(v.RWText)
	v.Unsubscribe(v.Responder)

	CopyFields(v, v2)

	if v.RWText == nil {
		v.RWText = text.New("")
	}
	v.Subscribe(v.RWText)

	if v.Responder == nil {
		v.Responder = &keyboard.Responder{}
	}
	v.Subscribe(v.Responder)
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

	st := text.NewStyledText(v.RWText.String(), style)

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

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Layouter:       &textInputLayouter{style: style, styledText: st, maxLines: v.MaxLines},
		Painter:        painter,
		NativeViewName: "gomatcha.io/matcha/view/textinput",
		NativeViewState: internal.MarshalProtobuf(&pbview.TextInput{
			Font:            style.Font().MarshalProtobuf(),
			StyledText:      st.MarshalProtobuf(),
			PlaceholderText: placeholderStyledText.MarshalProtobuf(),
			KeyboardType:    v.KeyboardType.MarshalProtobuf(),
			Focused:         v.Responder.Visible(),
			MaxLines:        int64(v.MaxLines),
			SecureTextEntry: v.Password,
		}),
		NativeFuncs: map[string]interface{}{
			"OnTextChange": func(data []byte) {
				pbevent := &pbview.TextInputEvent{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				v.RWText.UnmarshalProtobuf(pbevent.StyledText.Text)
				if v.OnChange != nil {
					v.OnChange(v.RWText)
				}
			},
			"OnSubmit": func() {
				if v.OnSubmit != nil {
					v.OnSubmit(v.RWText)
				}
			},
			"OnFocus": func(data []byte) {
				pbevent := &pbview.TextInputFocusEvent{}
				err := proto.Unmarshal(data, pbevent)
				if err != nil {
					fmt.Println("error", err)
					return
				}

				if pbevent.Focused {
					v.Responder.Show()
				} else {
					v.Responder.Dismiss()
				}
				if v.OnFocus != nil {
					v.OnFocus(v.Responder)
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
		st := text.NewStyledText("A", l.style)
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

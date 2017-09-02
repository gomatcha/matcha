package view

import (
	"image"
	"image/color"
	"runtime"
	"strings"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	pbbutton "gomatcha.io/matcha/pb/view/button"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
)

// Button implements a native button view.
type Button struct {
	Embed
	String     string
	Color      color.Color
	OnPress    func()
	Enabled    bool
	PaintStyle *paint.Style
}

// NewButton returns either the previous View in ctx with matching key, or a new View if none exists.
func NewButton() *Button {
	return &Button{
		Enabled: true,
	}
}

// Build implements view.View.
func (v *Button) Build(ctx *Context) Model {
	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Painter:        painter,
		Layouter:       &buttonLayouter{str: v.String},
		NativeViewName: "gomatcha.io/matcha/view/button",
		NativeViewState: &pbbutton.View{
			Str:     v.String,
			Enabled: v.Enabled,
			Color:   pb.ColorEncode(v.Color),
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

type buttonLayouter struct {
	str string
}

func (l *buttonLayouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	if runtime.GOOS == "android" {
		style := &text.Style{}
		style.SetFont(text.DefaultFont(14))
		st := text.NewStyledText(strings.ToUpper(l.str), style)
		size := st.Size(layout.Pt(0, 0), ctx.MaxSize, 1)

		const padding = 16.0
		g := layout.Guide{Frame: layout.Rt(0, 0, size.X+padding*2+16, 48)}
		return g, nil
	} else if runtime.GOOS == "darwin" {
		style := &text.Style{}
		style.SetFont(text.DefaultFont(20))
		st := text.NewStyledText(l.str, style)
		size := st.Size(layout.Pt(0, 0), ctx.MaxSize, 1)

		const padding = 10.0
		g := layout.Guide{Frame: layout.Rt(0, 0, size.X+padding*2, size.Y+padding*2)}
		return g, nil
	}
	return layout.Guide{}, nil
}

func (l *buttonLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *buttonLayouter) Unnotify(id comm.Id) {
	// no-op
}

type ImageButton struct {
	Embed
	Image      image.Image
	OnPress    func()
	PaintStyle *paint.Style
}

func NewImageButton() *ImageButton {
	return &ImageButton{}
}

func (v *ImageButton) Build(ctx *Context) Model {
	iv := NewImageView()
	iv.ResizeMode = ImageResizeModeCenter
	iv.Image = v.Image

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}

	t := &touch.ButtonRecognizer{
		OnTouch: func(e *touch.ButtonEvent) {
			if e.Kind == touch.EventKindRecognized && v.OnPress != nil {
				v.OnPress()
			}
		},
	}
	return Model{
		Children: []View{iv},
		Layouter: &imageButtonLayouter{},
		Painter:  painter,
		Options: []Option{
			touch.RecognizerList{t},
		},
	}
}

type imageButtonLayouter struct {
	str string
}

func (l *imageButtonLayouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := ctx.LayoutChild(0, ctx.MinSize, ctx.MaxSize)
	return g, []layout.Guide{g}
}

func (l *imageButtonLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *imageButtonLayouter) Unnotify(id comm.Id) {
	// no-op
}

// Package button implements a native button view.
package button

import (
	"image"
	"image/color"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	pbbutton "gomatcha.io/matcha/pb/view/button"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/imageview"
)

// View implements a native button view.
type View struct {
	view.Embed
	String     string
	Color      color.Color
	OnPress    func()
	Enabled    bool
	PaintStyle *paint.Style
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{
		Enabled: true,
		Color:   color.RGBA{14, 122, 254, 255},
	}
}

// Build implements view.View.
func (v *View) Build(ctx *view.Context) view.Model {
	style := &text.Style{}
	style.SetAlignment(text.AlignmentCenter)
	style.SetFont(text.Font{
		Name: "HelveticaNeue",
		Size: 20,
	})
	style.SetTextColor(v.Color)
	st := text.NewStyledText(v.String, style)

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
	styledText *text.StyledText
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

type ImageButton struct {
	view.Embed
	Image      image.Image
	OnPress    func()
	PaintStyle *paint.Style
}

func NewImageButton() *ImageButton {
	return &ImageButton{}
}

func (v *ImageButton) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	iv := imageview.New()
	iv.ResizeMode = imageview.ResizeModeCenter
	iv.Image = v.Image
	ivg := l.Add(iv, func(s *constraint.Solver) {
		s.Left(0)
		s.Top(0)
		s.HeightLess(l.MaxGuide().Height())
		s.WidthLess(l.MaxGuide().Width())
	})

	l.Solve(func(s *constraint.Solver) {
		s.WidthEqual(ivg.Width())
		s.HeightEqual(ivg.Height())
	})

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
	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  painter,
		Options: []view.Option{
			touch.RecognizerList{t},
		},
	}
}

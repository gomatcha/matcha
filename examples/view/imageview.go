// Package imageview provides examples of how to use the matcha/view packages.
package view

import (
	_ "image/jpeg"
	_ "image/png"

	"github.com/gomatcha/matcha/application"
	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/pointer"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view NewImageView", func() view.View {
		return NewImageView()
	})
}

type ImageView struct {
	view.Embed
	toggle bool
}

func NewImageView() *ImageView {
	return &ImageView{}
}

func (v *ImageView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := view.NewImageView()
	if v.toggle {
		chl.URL = "https://avatars0.githubusercontent.com/u/758035?v=4&s=460"
	} else {
		chl.URL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
	}
	chl.ResizeMode = view.ImageResizeModeFit
	chl.PaintStyle = &paint.Style{BackgroundColor: colornames.Pink}
	g1 := l.Add(chl, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(100)
		s.Width(200)
		s.Height(200)
	})
	_ = g1

	chl2 := view.NewImageView()
	chl2.Image = application.MustLoadImage("settings_airplane")
	chl2.ResizeMode = view.ImageResizeModeCenter
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
		s.WidthEqual(g1.Width())
		s.HeightEqual(g1.Height())
	})

	chl3 := view.NewImageView()
	chl3.Image = application.MustLoadImage("table_arrow")
	chl3.ResizeMode = view.ImageResizeModeStretch
	chl3.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
	chl3.ImageTint = colornames.Blue
	l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g2.Left())
		s.WidthEqual(g2.Width())
		s.HeightEqual(g2.Height())
	})

	tap := &pointer.TapGesture{
		Count: 1,
		OnEvent: func(e *pointer.TapEvent) {
			if e.Kind == pointer.EventKindRecognized {
				v.toggle = !v.toggle
				v.Signal()
			}
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Options: []view.Option{
			pointer.GestureList{tap},
		},
	}
}

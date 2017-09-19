// Package imageview provides examples of how to use the matcha/view packages.
package view

import (
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewImageView", func() view.View {
		return NewImageView()
	})
}

type ImageView struct {
	view.Embed
}

func NewImageView() *ImageView {
	return &ImageView{}
}

func (v *ImageView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := view.NewImageView()
	chl.URL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
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

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

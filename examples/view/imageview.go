// Package imageview provides examples of how to use the matcha/view packages.
package view

import (
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/imageview"
	"gomatcha.io/matcha/view/urlimageview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewImageView", func() *view.Root {
		return view.NewRoot(NewImageView(nil, ""))
	})
}

type ImageView struct {
	view.Embed
}

func NewImageView(ctx *view.Context, key string) *ImageView {
	if v, ok := ctx.Prev(key).(*ImageView); ok {
		return v
	}
	return &ImageView{
		Embed: view.Embed{Key: key},
	}
}

func (v *ImageView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := urlimageview.New()
	chl.URL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
	chl.ResizeMode = imageview.ResizeModeStretch
	g1 := l.Add(chl, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(100)
		s.Width(200)
		s.Height(200)
	})

	chl2 := imageview.New()
	chl2.Image = app.MustLoadImage("Airplane")
	chl2.ResizeMode = imageview.ResizeModeFit
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
		s.WidthEqual(g1.Width())
		s.HeightEqual(g1.Height())
	})

	chl3 := imageview.New()
	chl3.Image = app.MustLoadImage("TableArrow")
	chl3.ResizeMode = imageview.ResizeModeCenter
	chl3.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
	chl3.ImageTemplateColor = colornames.Blue
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

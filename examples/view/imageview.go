// Package imageview provides examples of how to use the matcha/view packages.
package view

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewImageView", func() view.View {
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

	label := view.NewTextView()
	label.String = "Tint Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	imageview := view.NewImageView()
	imageview.ImageTint = colornames.Red
	imageview.Image = application.MustLoadImage("checkbox_checked")
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Center Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview = view.NewImageView()
	imageview.Image = application.MustLoadImage("settings_airplane")
	imageview.ResizeMode = view.ImageResizeModeCenter
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Stretch Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview = view.NewImageView()
	imageview.Image = application.MustLoadImage("settings_airplane")
	imageview.ResizeMode = view.ImageResizeModeStretch
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Fill Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview = view.NewImageView()
	imageview.Image = application.MustLoadImage("settings_airplane")
	imageview.ResizeMode = view.ImageResizeModeFill
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Fit Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview = view.NewImageView()
	imageview.Image = application.MustLoadImage("settings_airplane")
	imageview.ResizeMode = view.ImageResizeModeFit
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Network Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.CenterXEqual(l.CenterX().Add(15))
	})

	imageview = view.NewImageView()
	imageview.URL = "https://avatars0.githubusercontent.com/u/758035?v=4&s=460"
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	label = view.NewTextView()
	label.String = "Toggle Image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview = view.NewImageView()
	if v.toggle {
		imageview.URL = "https://avatars1.githubusercontent.com/u/28016430?s=200&v=4"
	} else {
		imageview.URL = "https://avatars0.githubusercontent.com/u/758035?v=4&s=460"
	}
	imageview.ResizeMode = view.ImageResizeModeFill

	tap := &pointer.TapGesture{
		Count: 1,
		OnEvent: func(e *pointer.TapEvent) {
			fmt.Println("what")
			if e.Kind == pointer.EventKindRecognized {
				fmt.Println("recognize")
				v.toggle = !v.toggle
				v.Signal()
			}
		},
	}

	withOpts := view.WithOptions(imageview, pointer.GestureList{tap})
	g = l.Add(withOpts, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

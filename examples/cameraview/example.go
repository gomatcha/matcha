// Package custom provides examples of how to import custom components.
package customview

import (
	"image"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/customview NewView", func() view.View {
		return NewView()
	})
}

type View struct {
	view.Embed
	image       image.Image
	frontCamera bool
}

func NewView() *View {
	return &View{}
}

func (v *View) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Camera view:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	cam := NewCameraView()
	cam.FrontCamera = v.frontCamera
	cam.OnCapture = func(img image.Image) {
		v.image = img
		v.Signal()
	}
	g = l.Add(cam, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.Left(15)
		s.Right(-15)
		s.Height(200)
	})

	label = view.NewTextView()
	label.String = "Toggle front camera:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	switchview := view.NewSwitch()
	switchview.Value = v.frontCamera
	switchview.OnSubmit = func(value bool) {
		v.frontCamera = value
		v.Signal()
	}
	g = l.Add(switchview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Captured image:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	imageview := view.NewImageView()
	imageview.Image = v.image
	imageview.PaintStyle = &paint.Style{BackgroundColor: colornames.Black}
	g = l.Add(imageview, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.Left(15)
		s.Right(-15)
		s.Height(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

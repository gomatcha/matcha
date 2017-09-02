package view

import (
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	"gomatcha.io/matcha/pb/view/imageview"
)

type ImageResizeMode int

// TODO(KD): ResizeModeFit and ResizeModeFill behave oddly on Android.
const (
	// The image is resized proportionally such that a single axis is filled.
	ImageResizeModeFit ImageResizeMode = iota
	// The image is resized proportionally such that the entire view is filled.
	ImageResizeModeFill
	// The image is stretched to fill the view.
	ImageResizeModeStretch
	// The image is centered in the view with its natural size.
	ImageResizeModeCenter
)

func (m ImageResizeMode) MarshalProtobuf() imageview.ResizeMode {
	return imageview.ResizeMode(m)
}

// ImageView implements a view that displays an image.
type ImageView struct {
	Embed
	Image      image.Image
	ResizeMode ImageResizeMode
	ImageTint  color.Color
	PaintStyle *paint.Style
	image      image.Image
	pbImage    *pb.ImageOrResource
}

// NewImageView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewImageView() *ImageView {
	return &ImageView{}
}

// Build implements view.View.
func (v *ImageView) Build(ctx *Context) Model {
	if v.Image != v.image {
		v.image = v.Image
		v.pbImage = app.ImageMarshalProtobuf(v.image)
	}

	// Default to Center if we don't have an image
	bounds := image.Rect(0, 0, 0, 0)
	resizeMode := ImageResizeModeCenter
	scale := 1.0
	if v.image != nil {
		bounds = v.image.Bounds()
		resizeMode = v.ResizeMode

		if res, ok := v.image.(*app.ImageResource); ok {
			scale = res.Scale()
		}
	}

	var painter paint.Painter
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return Model{
		Painter:        painter,
		Layouter:       &imageViewLayouter{bounds: bounds, resizeMode: resizeMode, scale: scale},
		NativeViewName: "gomatcha.io/matcha/view/imageview",
		NativeViewState: &imageview.View{
			Image:      v.pbImage,
			Scale:      scale,
			ResizeMode: v.ResizeMode.MarshalProtobuf(),
			Tint:       pb.ColorEncode(v.ImageTint),
		},
	}
}

type imageViewLayouter struct {
	bounds     image.Rectangle
	scale      float64
	resizeMode ImageResizeMode
}

func (l *imageViewLayouter) Layout(ctx *layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{Frame: layout.Rect{Max: ctx.MinSize}}
	switch l.resizeMode {
	case ImageResizeModeFit:
		imgRatio := float64(l.bounds.Dx()) / l.scale / float64(l.bounds.Dy()) / l.scale
		maxRatio := ctx.MinSize.X / ctx.MinSize.Y
		if imgRatio > maxRatio {
			g.Frame.Max = layout.Pt(ctx.MinSize.X, ctx.MinSize.X/imgRatio)
		} else {
			g.Frame.Max = layout.Pt(ctx.MinSize.Y/imgRatio, ctx.MinSize.Y)
		}
	case ImageResizeModeFill:
		fallthrough
	case ImageResizeModeStretch:
		g.Frame.Max = ctx.MinSize
	case ImageResizeModeCenter:
		g.Frame.Max = layout.Pt(float64(l.bounds.Dx())/l.scale, float64(l.bounds.Dy())/l.scale)
	}
	return g, nil
}

func (l *imageViewLayouter) Notify(f func()) comm.Id {
	return 0 // no-op
}

func (l *imageViewLayouter) Unnotify(id comm.Id) {
	// no-op
}

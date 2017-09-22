package view

import (
	"context"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"gomatcha.io/matcha"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/paint"
	pb "gomatcha.io/matcha/proto"
	pbview "gomatcha.io/matcha/proto/view"
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

func (m ImageResizeMode) MarshalProtobuf() pbview.ImageResizeMode {
	return pbview.ImageResizeMode(m)
}

// ImageView implements a view that displays an image.
type ImageView struct {
	Embed
	Image      image.Image
	ResizeMode ImageResizeMode
	ImageTint  color.Color
	PaintStyle *paint.Style
	stage      Stage

	URL        string
	cancelFunc context.CancelFunc
	urlImage   image.Image
	err        error
	pbImage    *pb.ImageOrResource
}

// NewImageView returns either the previous View in ctx with matching key, or a new View if none exists.
func NewImageView() *ImageView {
	return &ImageView{}
}

// Build implements view.View.
func (v *ImageView) Build(ctx Context) Model {
	if v.pbImage == nil {
		if v.Image != nil {
			v.pbImage = internal.ImageMarshalProtobuf(v.Image)
		} else if v.urlImage != nil {
			v.pbImage = internal.ImageMarshalProtobuf(v.urlImage)
		}
	}

	// Default to Center if we don't have an image
	bounds := image.Rect(0, 0, 0, 0)
	resizeMode := ImageResizeModeCenter
	scale := 1.0
	if v.Image != nil {
		bounds = v.Image.Bounds()
		resizeMode = v.ResizeMode

		if res, ok := v.Image.(*application.ImageResource); ok {
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
		NativeViewState: internal.MarshalProtobuf(&pbview.ImageView{
			Image:      v.pbImage,
			Scale:      scale,
			ResizeMode: v.ResizeMode.MarshalProtobuf(),
			Tint:       pb.ColorEncode(v.ImageTint),
		}),
	}
}

// Lifecycle implements view.View.
func (v *ImageView) Lifecycle(from, to Stage) {
	v.stage = to
	v.reload()
}

func (v *ImageView) Update(v2 View) {
	if v2.(*ImageView).Image != v.Image {
		v.pbImage = nil
	}
	if v2.(*ImageView).URL != v.URL {
		v.cancel()
		v.urlImage = nil
		v.err = nil
	}
	CopyFields(v, v2)
}

func (v *ImageView) reload() {
	if v.stage < StageMounted {
		v.cancel()
		return
	}
	if v.URL == "" || v.Image != nil || v.cancelFunc != nil || v.urlImage != nil || v.err != nil {
		return
	}

	c, cancelFunc := context.WithCancel(context.Background())
	v.cancelFunc = cancelFunc
	go func(url string) {
		image, err := loadImageURL(url)

		matcha.MainLocker.Lock()
		defer matcha.MainLocker.Unlock()

		select {
		case <-c.Done():
		default:
			v.cancelFunc()
			v.cancelFunc = nil
			v.urlImage = image
			v.err = err
			v.Signal()
		}
	}(v.URL)
}

func (v *ImageView) cancel() {
	if v.cancelFunc != nil {
		v.cancelFunc()
		v.cancelFunc = nil
	}
}

type imageViewLayouter struct {
	bounds     image.Rectangle
	scale      float64
	resizeMode ImageResizeMode
}

func (l *imageViewLayouter) Layout(ctx layout.Context) (layout.Guide, []layout.Guide) {
	g := layout.Guide{Frame: layout.Rect{Max: ctx.MinSize()}}
	switch l.resizeMode {
	case ImageResizeModeFit:
		imgRatio := float64(l.bounds.Dx()) / l.scale / float64(l.bounds.Dy()) / l.scale
		maxRatio := ctx.MinSize().X / ctx.MinSize().Y
		if imgRatio > maxRatio {
			g.Frame.Max = layout.Pt(ctx.MinSize().X, ctx.MinSize().X/imgRatio)
		} else {
			g.Frame.Max = layout.Pt(ctx.MinSize().Y/imgRatio, ctx.MinSize().Y)
		}
	case ImageResizeModeFill:
		fallthrough
	case ImageResizeModeStretch:
		g.Frame.Max = ctx.MinSize()
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

func loadImageURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	return img, err
}

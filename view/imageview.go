package view

import (
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/gomatcha/matcha"
	"github.com/gomatcha/matcha/application"
	"github.com/gomatcha/matcha/comm"
	"github.com/gomatcha/matcha/internal"
	"github.com/gomatcha/matcha/layout"
	"github.com/gomatcha/matcha/paint"
	pb "github.com/gomatcha/matcha/proto"
	pbview "github.com/gomatcha/matcha/proto/view"
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
	URL        string
	ResizeMode ImageResizeMode
	ImageTint  color.Color
	PaintStyle *paint.Style

	cancelFunc context.CancelFunc
	err        error
	image      *pb.ImageOrResource
}

// NewImageView returns a new view.
func NewImageView() *ImageView {
	return &ImageView{}
}

// Lifecycle implements view.View.
func (v *ImageView) Lifecycle(from, to Stage) {
	if EntersStage(from, to, StageMounted) {
		v.begin()
	} else if ExitsStage(from, to, StageMounted) {
		v.end()
	}
}

func (v *ImageView) Update(v2 View) {
	prev := v2.(*ImageView)
	if prev.Image != v.Image || prev.URL != v.URL {
		v.end()
		CopyFields(v, v2)
		v.begin()
	} else {
		CopyFields(v, v2)
	}
}

// Build implements view.View.
func (v *ImageView) Build(ctx Context) Model {
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
		NativeViewName: "github.com/gomatcha/matcha/view/imageview",
		NativeViewState: internal.MarshalProtobuf(&pbview.ImageView{
			Image:      v.image,
			Scale:      scale,
			ResizeMode: v.ResizeMode.MarshalProtobuf(),
			Tint:       pb.ColorEncode(v.ImageTint),
		}),
	}
}

func (v *ImageView) begin() {
	if v.Image != nil {
		v.image = internal.ImageMarshalProtobuf(v.Image)
	} else if v.URL != "" {
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
				v.image = image
				v.err = err
				v.Signal()
			}
		}(v.URL)
	} else {
		v.err = errors.New("ImageView No Image or URL")
	}
}

func (v *ImageView) end() {
	if v.cancelFunc != nil {
		v.cancelFunc()
		v.cancelFunc = nil
	}
	v.image = nil
	v.err = nil
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

func loadImageURL(url string) (*pb.ImageOrResource, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("loadImageURL error", err)
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("decodeImage error", err)
	}
	return internal.ImageMarshalProtobuf(img), nil
}

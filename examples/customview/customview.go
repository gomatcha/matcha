package customview

import (
	"bytes"
	"fmt"
	"image"

	"golang.org/x/image/colornames"
	protocustomview "gomatcha.io/matcha/examples/customview/proto"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type CameraView struct {
	view.Embed
	FrontCamera bool
	OnCapture   func(img image.Image)
}

// NewCameraView returns an initialized CameraView instance.
func NewCameraView() *CameraView {
	return &CameraView{}
}

// Build implements view.View.
func (v *CameraView) Build(ctx view.Context) view.Model {
	return view.Model{
		Painter:        &paint.Style{BackgroundColor: colornames.Black},
		NativeViewName: "gomatcha.io/matcha/examples/customview CameraView",
		NativeViewState: internal.MarshalProtobuf(&protocustomview.View{
			FrontCamera: v.FrontCamera,
		}),
		NativeFuncs: map[string]interface{}{
			"OnCapture": func(data []byte) {
				buf := bytes.NewBuffer(data)
				img, _, err := image.Decode(buf)
				if err != nil {
					fmt.Println("error decoding image", err)
					return
				}

				if v.OnCapture != nil {
					v.OnCapture(img)
				}
			},
		},
	}
}

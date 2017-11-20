package cameraview

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

type nativeState struct {
	FrontCamera bool `json:"frontCamera"`
}

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
	state, _ := json.Marshal(&nativeState{FrontCamera: v.FrontCamera})

	return view.Model{
		Painter:         &paint.Style{BackgroundColor: colornames.Black},
		NativeViewName:  "gomatcha.io/matcha/examples/cameraview CameraView",
		NativeViewState: state,
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

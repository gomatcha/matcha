package application

import (
	"bytes"
	"image"
	"image/color"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/bridge"
	pb "gomatcha.io/matcha/proto"
	"gomatcha.io/matcha/proto/env"
)

// type Resource struct {
// 	path string
// }

// // Load loads a resource at path.
// func Load(path string) (*Resource, error) {
// 	return &Resource{path: path}, nil
// }

// // MustLoadImage loads the resource at path, or panics on error.
// func MustLoad(path string) *Resource {
// 	res, err := Load(path)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return res
// }

// func (r *Resource) MarshalProtobuf() *env.Resource {
// 	return &env.Resource{
// 		Path: r.path,
// 	}
// }

// ImageResource represents a static image asset and implements the image.Image interface.
type ImageResource struct {
	path  string
	rect  image.Rectangle
	image image.Image
	scale float64
}

// MustLoadImage loads the image at path.
func LoadImage(path string) (*ImageResource, error) {
	var propData []byte
	if runtime.GOOS == "android" {
		propData = bridge.Bridge("").Call("getPropertiesForResource", bridge.String(path)).ToInterface().([]byte)
	} else if runtime.GOOS == "darwin" {
		propData = bridge.Bridge("").Call("propertiesForResource:", bridge.String(path)).ToInterface().([]byte)
	}
	props := &pb.ImageProperties{}
	err := proto.Unmarshal(propData, props)
	if err != nil {
		return nil, err
	}

	return &ImageResource{
		path:  path,
		rect:  image.Rect(0, 0, int(props.Width), int(props.Height)),
		image: nil,
		scale: props.Scale,
	}, nil
}

// MustLoadImage loads the image at path, or panics on error.
func MustLoadImage(path string) *ImageResource {
	res, err := LoadImage(path)
	if err != nil {
		panic(err.Error())
	}
	return res
}

// ColorModel implements the image.Image interface.
func (res *ImageResource) ColorModel() color.Model {
	if res.image == nil {
		return color.RGBAModel
	}
	return res.image.ColorModel()
}

// Bounds implements the image.Image interface.
func (res *ImageResource) Bounds() image.Rectangle {
	return res.rect
}

// At implements the image.Image interface.
func (res *ImageResource) At(x, y int) color.Color {
	if res.image == nil {
		res.load()
	}
	return res.image.At(x, y)
}

// Scale returns the scale factor of the image.
func (res *ImageResource) Scale() float64 {
	return res.scale
}

// Path returns the path to the image.
func (res *ImageResource) Path() string {
	return res.path
}

func (res *ImageResource) load() {
	var data []byte
	if runtime.GOOS == "android" {
		data = bridge.Bridge("").Call("getImageForResource", bridge.String(res.path)).ToInterface().([]byte)
	} else if runtime.GOOS == "darwin" {
		data = bridge.Bridge("").Call("imageForResource:", bridge.String(res.path)).ToInterface().([]byte)
	}
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		res.image = image.NewUniform(colornames.Black)
		return
	}
	res.image = img
}

// MarshalProtobuf encodes res into a Protobuf object.
func (res *ImageResource) MarshalProtobuf() *env.ImageResource {
	if res == nil {
		return nil
	}
	return &env.ImageResource{
		Path: res.path,
	}
}

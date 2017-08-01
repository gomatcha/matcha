package app

import (
	"bytes"
	"image"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/gogo/protobuf/proto"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/pb"
	"gomatcha.io/matcha/pb/env"
)

type Resource struct {
	path string
}

// Load loads a resource at path.
func Load(path string) (*Resource, error) {
	return &Resource{path: path}, nil
}

// MustLoadImage loads the resource at path, or panics on error.
func MustLoad(path string) *Resource {
	res, err := Load(path)
	if err != nil {
		panic(err.Error())
	}
	return res
}

func (r *Resource) MarshalProtobuf() *env.Resource {
	return &env.Resource{
		Path: r.path,
	}
}

// Disable "Compress PNG Files" and "Remove Text Metadata from PNG Files" if loading image resources is not working.
type ImageResource struct {
	path  string
	rect  image.Rectangle
	image image.Image
	scale float64
}

// MustLoadImage loads the image at path.
func LoadImage(path string) (*ImageResource, error) {
	propData := bridge.Bridge().Call("propertiesForResource:", bridge.String(path)).ToInterface().([]byte)
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

func (res *ImageResource) load() {
	data := bridge.Bridge().Call("imageForResource:", bridge.String(res.path)).ToInterface().([]byte)
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		res.image = image.NewUniform(colornames.Black)
		return
	}
	res.image = img
}

func (res *ImageResource) MarshalProtobuf() *env.ImageResource {
	if res == nil {
		return nil
	}
	return &env.ImageResource{
		Path: res.path,
	}
}

func ImageMarshalProtobuf(img image.Image) *pb.ImageOrResource {
	if img == nil {
		return nil
	}
	if res, ok := img.(*ImageResource); ok {
		return &pb.ImageOrResource{
			Path: res.path,
		}
	} else {
		return &pb.ImageOrResource{
			Image: pb.ImageEncode(img),
		}
	}
}

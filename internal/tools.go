package internal

import (
	"os"
	"runtime/pprof"

	"image"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/pb"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/internal printStack", printStack)
}

func printStack() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func ImageMarshalProtobuf(img image.Image) *pb.ImageOrResource {
	if img == nil {
		return nil
	}
	if res, ok := img.(*application.ImageResource); ok {
		return &pb.ImageOrResource{
			Path: res.Path(),
		}
	} else {
		return &pb.ImageOrResource{
			Image: pb.ImageEncode(img),
		}
	}
}

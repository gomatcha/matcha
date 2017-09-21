package internal

import (
	"os"
	"runtime/pprof"

	"image"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/application"
	pb "gomatcha.io/matcha/proto"
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

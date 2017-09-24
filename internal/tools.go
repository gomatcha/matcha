package internal

import (
	"fmt"
	"os"
	"runtime/pprof"

	"image"

	gogoproto "github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/application"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/proto"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/internal printStack", printStack)
}

func printStack() {
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
}

func ImageMarshalProtobuf(img image.Image) *proto.ImageOrResource {
	if img == nil {
		return nil
	}
	if res, ok := img.(*application.ImageResource); ok {
		return &proto.ImageOrResource{
			Path: res.Path(),
		}
	} else {
		return &proto.ImageOrResource{
			Image: proto.ImageEncode(img),
		}
	}
}

func MarshalProtobuf(pb gogoproto.Message) []byte {
	data, err := gogoproto.Marshal(pb)
	if err != nil {
		fmt.Println("Error marshalling protobuf", pb, err)
	}
	return data
}

package view

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pbtouch "gomatcha.io/matcha/pb/touch"
	"gomatcha.io/matcha/touch"
)

type recognizer interface {
	marshalProtobuf() (proto.Message, map[string]interface{})
	equal(touch.Recognizer) bool
}

func init() {
	internal.RegisterMiddleware(func() interface{} { return &touchMiddleware{radix: radix.NewRadix()} })
}

type touchMiddleware struct {
	maxId int64
	radix *radix.Radix
}

func (r *touchMiddleware) MarshalProtobuf() proto.Message {
	return nil
}

func (r *touchMiddleware) Build(ctx *Context, next *Model) {
	path := idSliceToIntSlice(ctx.Path())
	node := r.radix.At(path)
	var prevIds map[int64]touch.Recognizer
	if node != nil {
		prevIds, _ = node.Value.(map[int64]touch.Recognizer)
	}

	ids := map[int64]touch.Recognizer{}

	var rs touch.RecognizerList
	for _, i := range next.Options {
		rs, _ = i.(touch.RecognizerList)
		if rs != nil {
			break
		}
	}

	// Diff prev and next recognizers
	for _, i := range rs {
		found := false
		for k, v := range prevIds {
			// Check that the id has not already been used.
			if _, ok := ids[k]; ok {
				continue
			}

			// Check that the recognizers are equal.
			if !i.(recognizer).equal(v) {
				continue
			}

			ids[k] = i
			found = true
		}

		// Generate a new id if we don't have a previous one.
		if !found {
			r.maxId += 1
			ids[r.maxId] = i
		}
	}

	if len(ids) == 0 {
		r.radix.Delete(path)
		return
	}

	// Add new list back to next.
	if node == nil {
		node = r.radix.Insert(path)
	}
	node.Value = ids

	// Serialize into protobuf.
	pbRecognizers := &pbtouch.RecognizerList{}
	allFuncs := map[string]interface{}{}
	for k, v := range ids {
		msg, funcs := v.(recognizer).marshalProtobuf()
		pbAny, err := ptypes.MarshalAny(msg)
		if err != nil {
			continue
		}

		pbRecognizer := &pbtouch.Recognizer{
			Id:         k,
			Recognizer: pbAny,
		}
		pbRecognizers.Recognizers = append(pbRecognizers.Recognizers, pbRecognizer)
		for k2, v2 := range funcs {
			allFuncs[k2] = v2
		}
	}

	if next.NativeValues == nil {
		next.NativeValues = map[string]proto.Message{}
	}
	next.NativeValues["gomatcha.io/matcha/touch"] = pbRecognizers

	if next.NativeFuncs == nil {
		next.NativeFuncs = map[string]interface{}{}
	}
	for k, v := range allFuncs {
		next.NativeFuncs[k] = v
	}
}

func (r *touchMiddleware) Key() string {
	return "gomatcha.io/matcha/touch"
}

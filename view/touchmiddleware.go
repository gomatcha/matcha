package view

import (
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pbtouch "gomatcha.io/matcha/pb/touch"
	"gomatcha.io/matcha/touch"
)

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

func (r *touchMiddleware) Build(ctx Context, next *Model) {
	path := idSliceToIntSlice(ctx.Path())
	node := r.radix.At(path)
	var prevIds map[int64]touch.Gesture
	if node != nil {
		prevIds, _ = node.Value.(map[int64]touch.Gesture)
	}

	ids := map[int64]touch.Gesture{}

	var rs touch.GestureList
	for _, i := range next.Options {
		rs, _ = i.(touch.GestureList)
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
			if internal.ReflectName(i) != internal.ReflectName(v) || i.TouchKey() != v.TouchKey() {
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
		model := v.Build()
		pbAny, err := ptypes.MarshalAny(model.NativeViewState)
		if err != nil {
			continue
		}

		pbRecognizer := &pbtouch.Recognizer{
			Id:         k,
			Recognizer: pbAny,
		}
		pbRecognizers.Recognizers = append(pbRecognizers.Recognizers, pbRecognizer)
		for k2, v2 := range model.NativeFuncs {
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

func idSliceToIntSlice(ids []Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

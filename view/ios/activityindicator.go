package ios

import (
	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pbapp "gomatcha.io/matcha/proto/app"
	"gomatcha.io/matcha/view"
)

// If any view has an ActivityIndicator option, the spinner will be visible.
//  return view.Model{
//      Options: []view.Option{app.ActivityIndicator{}}
//  }
type ActivityIndicator struct {
	// ActivityIndicator has no fields.
}

func (a *ActivityIndicator) OptionKey() string {
	return "gomatcha.io/matcha/app activity"
}

func init() {
	internal.RegisterMiddleware(func() interface{} {
		return &activityIndicatorMiddleware{
			radix: radix.NewRadix(),
		}
	})
}

type activityIndicatorMiddleware struct {
	radix *radix.Radix
}

func (m *activityIndicatorMiddleware) Build(ctx view.Context, model *view.Model) {
	path := idSliceToIntSlice(ctx.Path())

	add := false
	for _, i := range model.Options {
		if _, ok := i.(*ActivityIndicator); ok {
			add = true
		}
	}
	if add {
		m.radix.Insert(path)
	} else {
		m.radix.Delete(path)
	}
}

func (m *activityIndicatorMiddleware) MarshalProtobuf() proto.Message {
	visible := false
	m.radix.Range(func(path []int64, node *radix.Node) {
		visible = true
	})
	return &pbapp.ActivityIndicator{
		Visible: visible,
	}
}

func (m *activityIndicatorMiddleware) Key() string {
	return "gomatcha.io/matcha/app activity"
}

func idSliceToIntSlice(ids []view.Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

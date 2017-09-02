package view

import (
	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pbapp "gomatcha.io/matcha/pb/app"
)

func init() {
	internal.RegisterMiddleware(func() interface{} {
		return &statusBarMiddleware{
			radix: radix.NewRadix(),
		}
	})
}

type statusBarMiddleware struct {
	radix *radix.Radix
}

func (m *statusBarMiddleware) Build(ctx *Context, model *Model) {
	path := idSliceToIntSlice(ctx.Path())

	add := false
	statusBar := app.StatusBar{}
	for _, i := range model.Options {
		var ok bool
		if statusBar, ok = i.(app.StatusBar); ok {
			add = true
		}
	}
	if add {
		n := m.radix.Insert(path)
		n.Value = statusBar
	} else {
		m.radix.Delete(path)
	}
}

func (m *statusBarMiddleware) MarshalProtobuf() proto.Message {
	var statusBar app.StatusBar
	maxId := int64(-1)
	m.radix.Range(func(path []int64, node *radix.Node) {
		if len(path) > 0 && path[len(path)-1] > maxId {
			maxId = path[len(path)-1]
			statusBar, _ = node.Value.(app.StatusBar)
		}
	})
	return &pbapp.StatusBar{
		Hidden: statusBar.Hidden,
		Style:  pbapp.StatusBarStyle(statusBar.Style),
	}
}

func (m *statusBarMiddleware) Key() string {
	return "gomatcha.io/matcha/app statusbar"
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

func (m *activityIndicatorMiddleware) Build(ctx *Context, model *Model) {
	path := idSliceToIntSlice(ctx.Path())

	add := false
	for _, i := range model.Options {
		if _, ok := i.(app.ActivityIndicator); ok {
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

func idSliceToIntSlice(ids []Id) []int64 {
	ints := make([]int64, len(ids))
	for idx, i := range ids {
		ints[idx] = int64(i)
	}
	return ints
}

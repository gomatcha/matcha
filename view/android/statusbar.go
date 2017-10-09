package android

import (
	"fmt"
	"image/color"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pb "gomatcha.io/matcha/proto"
	pbapp "gomatcha.io/matcha/proto/view/android"
	"gomatcha.io/matcha/view"
)

type StatusBarStyle int

const (
	// Statusbar with light icons
	StatusBarStyleLight StatusBarStyle = iota
	// Statusbar with dark icons
	StatusBarStyleDark
)

// If multiple views have a statusBar, the most recently mounted one will be used.
//  return view.Model{
//      Options: []view.Option{
//          &android.StatusBar{ Color: colornames.Red },
//      },
//  }
type StatusBar struct {
	Style StatusBarStyle
	Color color.Color
}

func (s *StatusBar) OptionKey() string {
	return "gomatcha.io/matcha/view/android statusbar"
}

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

func (m *statusBarMiddleware) Build(ctx view.Context, model *view.Model) {
	path := idSliceToIntSlice(ctx.Path())

	var statusBar *StatusBar
	if model != nil {
		for _, i := range model.Options {
			if bar, ok := i.(*StatusBar); ok {
				statusBar = bar
			}
		}
	}
	if statusBar != nil {
		n := m.radix.Insert(path)
		n.Value = statusBar
	} else {
		m.radix.Delete(path)
	}
}

func (m *statusBarMiddleware) MarshalProtobuf() proto.Message {
	statusBar := &StatusBar{Color: colornames.Black, Style: StatusBarStyleLight}
	maxId := int64(-1)
	m.radix.Range(func(path []int64, node *radix.Node) {
		if len(path) > 0 && path[len(path)-1] > maxId {
			maxId = path[len(path)-1]
			statusBar = node.Value.(*StatusBar)
		}
	})
	return &pbapp.StatusBar{
		Style: statusBar.Style == StatusBarStyleLight,
		Color: pb.ColorEncode(statusBar.Color),
	}
}

func (m *statusBarMiddleware) Key() string {
	return "gomatcha.io/matcha/view/android statusbar"
}

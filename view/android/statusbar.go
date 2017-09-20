package android

import (
	"image/color"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	"gomatcha.io/matcha/pb"
	pbapp "gomatcha.io/matcha/pb/view/android"
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
// UIViewControllerBasedStatusBarAppearance must be set to False in the app's Info.plist
// to use this component.
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

	add := false
	statusBar := &StatusBar{}
	for _, i := range model.Options {
		var ok bool
		if statusBar, ok = i.(*StatusBar); ok {
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
	statusBar := &StatusBar{Color: colornames.Black}
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

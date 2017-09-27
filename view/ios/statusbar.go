package ios

import (
	"github.com/gogo/protobuf/proto"
	"gomatcha.io/matcha/internal"
	"gomatcha.io/matcha/internal/radix"
	pbapp "gomatcha.io/matcha/proto/app"
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
//          &ios.StatusBar{ Style: ios.StatusBarStyleLight },
//      },
//  }
type StatusBar struct {
	Hidden bool
	Style  StatusBarStyle
}

func (s *StatusBar) OptionKey() string {
	return "gomatcha.io/matcha/app statusbar"
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
	statusBar := &StatusBar{Style: StatusBarStyleDark}
	maxId := int64(-1)
	m.radix.Range(func(path []int64, node *radix.Node) {
		if len(path) > 0 && path[len(path)-1] > maxId {
			maxId = path[len(path)-1]
			statusBar = node.Value.(*StatusBar)
		}
	})
	if statusBar == nil {
		statusBar = &StatusBar{Style: StatusBarStyleDark}
	}
	return &pbapp.StatusBar{
		Hidden: statusBar.Hidden,
		Style:  pbapp.StatusBarStyle(statusBar.Style + 1),
	}
}

func (m *statusBarMiddleware) Key() string {
	return "gomatcha.io/matcha/app statusbar"
}

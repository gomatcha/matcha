package app

// func init() {
// 	internal.RegisterMiddleware(func() interface{} {
// 		return &statusBarMiddleware{
// 			radix: radix.NewRadix(),
// 		}
// 	})
// }

type StatusBarStyle int

const (
	StatusBarStyleDefault StatusBarStyle = iota
	StatusBarStyleLight
	StatusBarStyleDark
)

// If multiple views have a statusBar, the most recently mounted one will be used.
// UIViewControllerBasedStatusBarAppearance must be set to False in the app's Info.plist
// to use this component.
//  return view.Model{
//  	Options: []view.Option{
//  		app.StatusBar{ Style: app.StatusBarStyleLight },
//  	},
//  }
type StatusBar struct {
	Hidden bool
	Style  StatusBarStyle
}

func (s StatusBar) OptionKey() string {
	return "gomatcha.io/matcha/app statusbar"
}

// type statusBarMiddleware struct {
// 	radix *radix.Radix
// }

// func (m *statusBarMiddleware) Build(ctx *view.Context, model *view.Model) {
// 	path := idSliceToIntSlice(ctx.Path())

// 	add := false
// 	statusBar := StatusBar{}
// 	for _, i := range model.Options {
// 		var ok bool
// 		if statusBar, ok = i.(StatusBar); ok {
// 			add = true
// 		}
// 	}
// 	if add {
// 		n := m.radix.Insert(path)
// 		n.Value = statusBar
// 	} else {
// 		m.radix.Delete(path)
// 	}
// }

// func (m *statusBarMiddleware) MarshalProtobuf() proto.Message {
// 	var statusBar StatusBar
// 	maxId := int64(-1)
// 	m.radix.Range(func(path []int64, node *radix.Node) {
// 		if len(path) > 0 && path[len(path)-1] > maxId {
// 			maxId = path[len(path)-1]
// 			statusBar, _ = node.Value.(StatusBar)
// 		}
// 	})
// 	return &app.StatusBar{
// 		Hidden: statusBar.Hidden,
// 		Style:  app.StatusBarStyle(statusBar.Style),
// 	}
// }

// func (m *statusBarMiddleware) Key() string {
// 	return "gomatcha.io/matcha/app statusbar"
// }

// func init() {
// 	internal.RegisterMiddleware(func() interface{} {
// 		return &activityIndicatorMiddleware{
// 			radix: radix.NewRadix(),
// 		}
// 	})
// }

// If any view has an ActivityIndicator option, the spinner will be visible.
//  return view.Model{
//  	Options: []view.Option{app.ActivityIndicator{}}
//  }
type ActivityIndicator struct {
	// ActivityIndicator has no fields.
}

func (a ActivityIndicator) OptionKey() string {
	return "gomatcha.io/matcha/app activity"
}

// type activityIndicatorMiddleware struct {
// 	radix *radix.Radix
// }

// func (m *activityIndicatorMiddleware) Build(ctx *view.Context, model *view.Model) {
// 	path := idSliceToIntSlice(ctx.Path())

// 	add := false
// 	for _, i := range model.Options {
// 		if _, ok := i.(ActivityIndicator); ok {
// 			add = true
// 		}
// 	}
// 	if add {
// 		m.radix.Insert(path)
// 	} else {
// 		m.radix.Delete(path)
// 	}
// }

// func (m *activityIndicatorMiddleware) MarshalProtobuf() proto.Message {
// 	visible := false
// 	m.radix.Range(func(path []int64, node *radix.Node) {
// 		visible = true
// 	})
// 	return &app.ActivityIndicator{
// 		Visible: visible,
// 	}
// }

// func (m *activityIndicatorMiddleware) Key() string {
// 	return "gomatcha.io/matcha/app activity"
// }

// func idSliceToIntSlice(ids []view.Id) []int64 {
// 	ints := make([]int64, len(ids))
// 	for idx, i := range ids {
// 		ints[idx] = int64(i)
// 	}
// 	return ints
// }

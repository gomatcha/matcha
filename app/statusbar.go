package app

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

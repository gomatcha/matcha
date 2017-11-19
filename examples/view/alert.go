package view

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewAlertView", func() view.View {
		return NewAlertView()
	})
}

type AlertView struct {
	view.Embed
	showView bool
}

func NewAlertView() *AlertView {
	return &AlertView{}
}

func (v *AlertView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Alert:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	button := view.NewButton()
	button.String = "Show"
	button.OnPress = func() {
		view.Alert("Title", "")
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Alert w/ Message:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Show"
	button.OnPress = func() {
		view.Alert("Title", "Message")
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Alert w/ Button:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Show"
	button.OnPress = func() {
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress") }}
		view.Alert("Title", "Message", ok)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Alert w/ 2 Buttons:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Show"
	button.OnPress = func() {
		cancel := &view.AlertButton{Title: "Cancel", OnPress: func() { fmt.Println("OnPress Cancel") }}
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress OK") }}
		view.Alert("Title", "Message", ok, cancel)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Alert w/ 3 Buttons:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button = view.NewButton()
	button.String = "Show"
	button.OnPress = func() {
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress OK") }}
		other := &view.AlertButton{Title: "Other", OnPress: func() { fmt.Println("OnPress Other") }}
		cancel := &view.AlertButton{Title: "Cancel", OnPress: func() { fmt.Println("OnPress Cancel") }}
		view.Alert("Title", "Message", ok, cancel, other)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

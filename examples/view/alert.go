package view

import (
	"fmt"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view NewAlertView", func() view.View {
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

	chl1 := view.NewButton()
	chl1.String = "Alert"
	chl1.OnPress = func() {
		view.Alert("Title", "")
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(0)
		s.Width(200)
	})

	chl2 := view.NewButton()
	chl2.String = "Alert with Message"
	chl2.OnPress = func() {
		view.Alert("Title", "Message")
	}
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.Left(0)
		s.Width(200)
	})

	chl3 := view.NewButton()
	chl3.String = "Alert with Button"
	chl3.OnPress = func() {
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress") }}
		view.Alert("Title", "Message", ok)
	}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.Left(0)
		s.Width(200)
	})

	chl4 := view.NewButton()
	chl4.String = "Alert with 2 Button"
	chl4.OnPress = func() {
		cancel := &view.AlertButton{Title: "Cancel", OnPress: func() { fmt.Println("OnPress Cancel") }}
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress OK") }}
		view.Alert("Title", "Message", ok, cancel)
	}
	g4 := l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g3.Bottom())
		s.Left(0)
		s.Width(200)
	})

	chl5 := view.NewButton()
	chl5.String = "Alert with 3 Button"
	chl5.OnPress = func() {
		ok := &view.AlertButton{Title: "OK", OnPress: func() { fmt.Println("OnPress OK") }}
		other := &view.AlertButton{Title: "Other", OnPress: func() { fmt.Println("OnPress Other") }}
		cancel := &view.AlertButton{Title: "Cancel", OnPress: func() { fmt.Println("OnPress Cancel") }}
		view.Alert("Title", "Message", ok, cancel, other)
	}
	_ = l.Add(chl5, func(s *constraint.Solver) {
		s.TopEqual(g4.Bottom())
		s.Left(0)
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

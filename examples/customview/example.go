// Package custom provides examples of how to import custom components.
package customview

import (
	"fmt"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/customview NewView", func() view.View {
		return NewView()
	})
}

type View struct {
	view.Embed
}

func NewView() *View {
	return &View{}
}

func (v *View) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	chl1 := NewCustomView()
	chl1.OnSubmit = func(v bool) {
		fmt.Println("OnSubmit", v)
	}
	l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

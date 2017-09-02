// Package table provides examples of how to use the matcha/layout/table package.
package table

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/scrollview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/table New", func() *view.Root {
		return view.NewRoot(New())
	})
}

type TableView struct {
	view.Embed
}

func New() *TableView {
	return &TableView{}
}

func (v *TableView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	childLayouter := &table.Layouter{}
	for i := 0; i < 20; i++ {
		childView := NewTableCell()
		childView.String = "TEST TEST"
		childView.Painter = &paint.Style{BackgroundColor: colornames.Red}
		childLayouter.Add(childView, nil)
	}

	scrollView := scrollview.New()
	scrollView.PaintStyle = &paint.Style{BackgroundColor: colornames.Cyan}
	scrollView.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	scrollView.ContentLayouter = childLayouter
	scrollView.ContentChildren = childLayouter.Views()
	_ = l.Add(scrollView, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(0))
		s.LeftEqual(constraint.Const(0))
		s.WidthEqual(constraint.Const(400))
		s.HeightEqual(constraint.Const(400))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}

}

type TableCell struct {
	view.Embed
	String  string
	Painter paint.Painter
}

func NewTableCell() *TableCell {
	return &TableCell{}
}

func (v *TableCell) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.HeightEqual(constraint.Const(50))
	})

	textView := view.NewTextView()
	textView.String = v.String
	textView.Style.SetFont(text.Font{
		Name: "HelveticaNeue",
		Size: 20,
	})
	l.Add(textView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.CenterYEqual(l.CenterY())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  v.Painter,
	}
}

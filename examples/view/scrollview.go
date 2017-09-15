package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewScrollView", func() view.View {
		return NewScrollView()
	})
}

type ScrollView struct {
	view.Embed
}

func NewScrollView() *ScrollView {
	return &ScrollView{}
}

func (v *ScrollView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	childLayouter := &table.Layouter{}
	for i := 0; i < 20; i++ {
		childLayouter.Add(NewTableCell(), nil)
	}

	scrollview := view.NewScrollView()
	scrollview.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
	scrollview.ContentLayouter = childLayouter
	scrollview.ContentChildren = childLayouter.Views()
	_ = l.Add(scrollview, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(300)
		s.Height(500)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type TableCell struct {
	view.Embed
}

func NewTableCell() *TableCell {
	return &TableCell{}
}

func (v *TableCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
	})

	chl := view.NewBasicView()
	chl.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	l.Add(chl, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.TopEqual(l.Top().Add(10))
		s.BottomEqual(l.Bottom().Add(-10))
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

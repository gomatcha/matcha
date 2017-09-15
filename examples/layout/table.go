package layout

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/layout NewTableView", func() view.View {
		return NewTableView()
	})
}

type TableView struct {
	view.Embed
}

func NewTableView() *TableView {
	return &TableView{}
}

func (v *TableView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	childLayouter := &table.Layouter{
		StartEdge: layout.EdgeLeft,
	}
	for i := 0; i < 20; i++ {
		childView := NewTableCell()
		childView.String = fmt.Sprintf("Cell %v", i)
		childView.Painter = &paint.Style{BackgroundColor: colornames.Red}
		childLayouter.Add(childView, nil)
	}

	sv := view.NewScrollView()
	sv.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	sv.ContentLayouter = childLayouter
	sv.ContentChildren = childLayouter.Views()
	sv.ScrollAxes = layout.AxisX
	sv.PaintStyle = &paint.Style{BackgroundColor: colornames.Cyan}
	_ = l.Add(sv, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.WidthEqual(l.Width())
		s.HeightEqual(l.Height())
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

func (v *TableCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
		s.Width(100)
	})

	textView := view.NewTextView()
	textView.String = v.String
	textView.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
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

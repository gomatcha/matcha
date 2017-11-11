package layout

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
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

	verticalLayouter := &table.Layouter{StartEdge: layout.EdgeBottom}
	for i := 0; i < 20; i++ {
		childView := NewTableCell()
		childView.String = fmt.Sprintf("Cell %v", i)
		verticalLayouter.Add(childView, nil)
	}
	verticalSV := view.NewScrollView()
	verticalSV.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	verticalSV.ContentLayouter = verticalLayouter
	verticalSV.ContentChildren = verticalLayouter.Views()
	verticalSV.ScrollAxes = layout.AxisY
	g := l.Add(verticalSV, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.WidthEqual(l.Width())
		s.Height(400)
	})

	horizontalLayouter := &table.Layouter{StartEdge: layout.EdgeLeft}
	for i := 0; i < 20; i++ {
		childView := NewTableCell()
		childView.String = fmt.Sprintf("Cell %v", i)
		horizontalLayouter.Add(childView, nil)
	}
	horizontalSV := view.NewScrollView()
	horizontalSV.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	horizontalSV.ContentLayouter = horizontalLayouter
	horizontalSV.ContentChildren = horizontalLayouter.Views()
	horizontalSV.ScrollAxes = layout.AxisX
	_ = l.Add(horizontalSV, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.Left(0)
		s.WidthEqual(l.Width())
		s.BottomEqual(l.Bottom())
	})

	return view.Model{
		Children: l.Views(),
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Layouter: l,
	}

}

type TableCell struct {
	view.Embed
	String string
}

func NewTableCell() *TableCell {
	return &TableCell{}
}

func (v *TableCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.HeightGreater(constraint.Const(100))
		s.WidthGreater(constraint.Const(100))
	})

	border := view.NewBasicView()
	border.Painter = &paint.Style{BackgroundColor: colornames.Green}
	l.Add(border, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.TopEqual(l.Top().Add(10))
		s.BottomEqual(l.Bottom().Add(-10))
	})

	textView := view.NewTextView()
	textView.String = v.String
	textView.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	textView.Style.SetAlignment(text.AlignmentCenter)
	l.Add(textView, func(s *constraint.Solver) {
		s.LeftEqual(l.Left().Add(10))
		s.RightEqual(l.Right().Add(-10))
		s.CenterYEqual(l.CenterY())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

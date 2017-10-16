package view

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout"
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
	scrollPosition *view.ScrollPosition
}

func NewScrollView() *ScrollView {
	return &ScrollView{
		scrollPosition: &view.ScrollPosition{},
	}
}

func (v *ScrollView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageMounted) {
		v.Subscribe(v.scrollPosition)
	} else if view.ExitsStage(from, to, view.StageMounted) {
		v.Unsubscribe(v.scrollPosition)
	}
}

func (v *ScrollView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	vtable := &table.Layouter{}
	for i := 0; i < 5; i++ {
		cell := NewTableCell()
		cell.Axis = layout.AxisY
		cell.Index = i
		vtable.Add(cell, nil)
	}

	scrollview := view.NewScrollView()
	scrollview.ScrollPosition = v.scrollPosition
	scrollview.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
	scrollview.ContentLayouter = vtable
	scrollview.ContentChildren = vtable.Views()
	g1 := l.Add(scrollview, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(200)
		s.BottomEqual(l.Bottom().Add(-210))
	})

	htable := &table.Layouter{
		StartEdge: layout.EdgeLeft,
	}
	for i := 0; i < 5; i++ {
		cell := NewTableCell()
		cell.Axis = layout.AxisX
		cell.Index = i
		htable.Add(cell, nil)
	}

	hscrollview := view.NewScrollView()
	hscrollview.ScrollAxes = layout.AxisX
	hscrollview.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
	hscrollview.ContentLayouter = htable
	hscrollview.ContentChildren = htable.Views()
	_ = l.Add(hscrollview, func(s *constraint.Solver) {
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.Height(200)
		s.BottomEqual(l.Bottom())
	})

	textView := view.NewTextView()
	textView.PaintStyle = &paint.Style{BackgroundColor: colornames.Red}
	textView.String = fmt.Sprintln("Position:", v.scrollPosition.X.Value(), v.scrollPosition.Y.Value())
	textView.MaxLines = 2
	g3 := l.Add(textView, func(s *constraint.Solver) {
		s.Top(50)
		s.LeftEqual(g1.Right())
		s.RightEqual(l.Right())
		s.Height(100)
	})

	button := view.NewButton()
	button.String = "Scroll"
	button.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	button.OnPress = func() {
		fmt.Println("OnPress")
		a := &animate.Basic{
			Start: v.scrollPosition.Y.Value(),
			End:   200,
			Dur:   time.Second / 5,
		}
		v.scrollPosition.Y.Run(a)
	}
	_ = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g3.Bottom())
		s.LeftEqual(g1.Right())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type TableCell struct {
	view.Embed
	Axis  layout.Axis
	Index int
}

func NewTableCell() *TableCell {
	return &TableCell{}
}

func (v *TableCell) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		if v.Axis == layout.AxisY {
			s.Height(200)
		} else {
			s.Width(200)
		}
	})

	label := view.NewTextView()
	label.String = strconv.Itoa(v.Index)
	l.Add(label, func(s *constraint.Solver) {
	})

	border := view.NewBasicView()
	border.Painter = &paint.Style{BackgroundColor: colornames.Gray}
	l.Add(border, func(s *constraint.Solver) {
		s.Height(1)
		s.LeftEqual(l.Left())
		s.RightEqual(l.Right())
		s.BottomEqual(l.Bottom())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

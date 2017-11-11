package view

import (
	"strconv"
	"time"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
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

	label := view.NewTextView()
	label.String = "Position: " + strconv.Itoa(int(v.scrollPosition.X.Value())) + ", " + strconv.Itoa(int(v.scrollPosition.Y.Value()))
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	button := view.NewButton()
	button.String = "Scroll"
	button.OnPress = func() {
		a := &animate.Basic{
			Start: v.scrollPosition.Y.Value(),
			End:   500,
			Dur:   time.Second / 5,
		}
		v.scrollPosition.Y.Run(a)
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	table := &table.Layouter{}
	for i := 0; i < 20; i++ {
		cell := NewTableCell()
		cell.String = strconv.Itoa(i)
		table.Add(cell, nil)
	}

	sv := view.NewScrollView()
	sv.ScrollPosition = v.scrollPosition
	sv.ContentLayouter = table
	sv.ContentChildren = table.Views()
	g = l.Add(sv, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(200)
		s.RightEqual(l.Right())
		s.BottomEqual(l.Bottom())
	})

	// textView := view.NewTextView()
	// textView.String = fmt.Sprintln("Scroll Position:", int(v.scrollPosition.X.Value()), int(v.scrollPosition.Y.Value()))
	// g3 := l.Add(textView, func(s *constraint.Solver) {
	// 	s.Top(50)
	// 	s.LeftEqual(g1.Right())
	// 	s.RightEqual(l.Right())
	// })

	// button := view.NewButton()
	// button.String = "Scroll to point"
	// button.OnPress = func() {
	// 	a := &animate.Basic{
	// 		Start: v.scrollPosition.Y.Value(),
	// 		End:   500,
	// 		Dur:   time.Second / 5,
	// 	}
	// 	v.scrollPosition.Y.Run(a)
	// }
	// _ = l.Add(button, func(s *constraint.Solver) {
	// 	s.TopEqual(g3.Bottom().Add(10))
	// 	s.LeftEqual(g1.Right())
	// 	s.RightEqual(l.Right())
	// })

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
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

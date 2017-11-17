package ios

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view/ios NewSegmentView", func() view.View {
		return NewSegmentView()
	})
}

type SegmentView struct {
	view.Embed
	value int
}

func NewSegmentView() *SegmentView {
	return &SegmentView{
		value: 1,
	}
}

func (v *SegmentView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Segmented view:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
	})

	segmented := ios.NewSegmentView()
	segmented.Titles = []string{"title1", "title2", "title3"}
	segmented.Value = v.value
	segmented.OnChange = func(value int) {
		v.value = value
		v.Signal()
		fmt.Println("onChange", value)
	}
	g = l.Add(segmented, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	label = view.NewTextView()
	label.String = "Momentary segmented view:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	segmented = ios.NewSegmentView()
	segmented.Titles = []string{"title1", "title2", "title3"}
	segmented.Value = v.value
	segmented.Momentary = true
	segmented.OnChange = func(value int) {
		v.value = value
		v.Signal()
		fmt.Println("onChange", value)
	}
	g = l.Add(segmented, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	label = view.NewTextView()
	label.String = "Disabled segmented view:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	segmented = ios.NewSegmentView()
	segmented.Titles = []string{"title1", "title2", "title3"}
	segmented.Value = v.value
	segmented.Enabled = false
	l.Add(segmented, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(200)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

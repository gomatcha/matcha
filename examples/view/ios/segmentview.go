package ios

import (
	"fmt"

	"github.com/gomatcha/matcha/bridge"
	"github.com/gomatcha/matcha/layout/constraint"
	"github.com/gomatcha/matcha/paint"
	"github.com/gomatcha/matcha/view"
	"github.com/gomatcha/matcha/view/ios"
	"golang.org/x/image/colornames"
)

func init() {
	bridge.RegisterFunc("github.com/gomatcha/matcha/examples/view/ios NewSegmentView", func() view.View {
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

	chl1 := ios.NewSegmentView()
	chl1.Titles = []string{"title1", "title2", "title3"}
	chl1.Value = v.value
	chl1.Momentary = true
	chl1.OnChange = func(value int) {
		v.value = value
		v.Signal()
		fmt.Println("onChange", value)
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
	})

	chl2 := ios.NewSegmentView()
	chl2.Titles = []string{"title1", "title2", "title3"}
	chl2.Value = v.value
	chl2.Enabled = false
	l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom().Add(50))
		s.LeftEqual(g1.Left())
		s.WidthEqual(g1.Width())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

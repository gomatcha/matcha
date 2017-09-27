// Package complex provides examples of many matcha subpackages working together.
package complex

import (
	"fmt"
	"time"

	"golang.org/x/image/colornames"

	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/ios"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/complex NewNestedView", func() view.View {
		return NewNestedView()
	})
}

type NestedView struct {
	view.Embed
	counter      int
	sliderValue  comm.Float64Value
	segmentValue int
	value        animate.Value
}

func NewNestedView() *NestedView {
	return &NestedView{}
}

func (v *NestedView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageVisible) {
		v.value.Run(&animate.Basic{
			Start: 0,
			End:   1,
			Ease:  animate.DefaultEase,
			Dur:   2 * time.Second,
		})
	}
}

func (v *NestedView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	value := animate.FloatLerp{Start: 0, End: 150}.Notifier(&v.value)
	// value := animate.FloatInterpolate(animate.FloatLerp{Start: 0, End: 150}, &v.value)

	chl1 := view.NewBasicView()
	chl1.Painter = &paint.AnimatedStyle{
		BackgroundColor: animate.RGBALerp{Start: colornames.Red, End: colornames.Yellow}.Notifier(&v.value),
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.WidthEqual(constraint.Notifier(value))
		s.HeightEqual(constraint.Notifier(value))
	})

	chl2 := view.NewBasicView()
	chl2.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
		s.Width(300)
		s.Height(300)
	})

	chl3 := view.NewBasicView()
	chl3.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g2.Left())
		s.Width(100)
		s.Height(100)
	})

	chl4 := view.NewBasicView()
	chl4.Painter = &paint.Style{BackgroundColor: colornames.Magenta}
	g4 := l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g3.Right())
		s.Width(50)
		s.Height(50)
	})

	chl5 := view.NewTextView()
	chl5.String = "Subtitle"
	chl5.Style.SetAlignment(text.AlignmentCenter)
	chl5.Style.SetStrikethroughStyle(text.StrikethroughStyleSingle)
	chl5.Style.SetStrikethroughColor(colornames.Magenta)
	chl5.Style.SetUnderlineStyle(text.UnderlineStyleDouble)
	chl5.Style.SetUnderlineColor(colornames.Green)
	chl5.Style.SetFont(text.FontWithName("AmericanTypewriter-Bold", 20))
	chl5p := view.WithPainter(chl5, &paint.Style{BackgroundColor: colornames.Cyan})

	g5 := l.Add(chl5p, func(s *constraint.Solver) {
		s.BottomEqual(g2.Bottom())
		s.RightEqual(g2.Right().Add(-15))
	})

	chl6 := view.NewTextView()
	chl6.String = fmt.Sprintf("Counter: %v", v.counter)
	chl6.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	chl6p := view.WithPainter(chl6, &paint.Style{BackgroundColor: colornames.Red})
	g6 := l.Add(chl6p, func(s *constraint.Solver) {
		s.BottomEqual(g5.Top())
		s.RightEqual(g2.Right().Add(-15))
	})

	chl8 := view.NewButton()
	chl8.String = "Button"
	chl8.OnPress = func() {
		v.counter += 1
		v.Signal()

		view.Alert("Alert", "Message",
			&view.AlertButton{
				Title: "OK",
				OnPress: func() {
					fmt.Println("OK")
				},
			},
			&view.AlertButton{
				Title: "Cancel",
				OnPress: func() {
					fmt.Println("Cancel")
				},
			},
		)
	}
	g8 := l.Add(chl8, func(s *constraint.Solver) {
		s.BottomEqual(g6.Top())
		s.RightEqual(g2.Right().Add(-15))
	})

	if v.counter%2 == 0 {
		chl9 := view.NewImageView()
		chl9.URL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
		chl9.ResizeMode = view.ImageResizeModeFit
		pChl9 := view.WithPainter(chl9, &paint.Style{BackgroundColor: colornames.Cyan})

		_ = l.Add(pChl9, func(s *constraint.Solver) {
			s.BottomEqual(g8.Top())
			s.RightEqual(g2.Right().Add(-15))
			s.Width(200)
			s.Height(200)
		})
	}
	chl11 := view.NewSwitch()
	chl11.OnSubmit = func(value bool) {
		a := 0.0
		if value {
			a = 1.0
		}

		v.sliderValue.SetValue(0.2)
		v.value.Run(&animate.Basic{
			Start: v.value.Value(),
			End:   a,
			Ease:  animate.DefaultEase,
			Dur:   2 * time.Second,
		})
	}
	_ = l.Add(chl11, func(s *constraint.Solver) {
		s.LeftEqual(g6.Right())
		s.TopEqual(g6.Top())
	})

	childLayouter := &table.Layouter{}
	childLayouter.StartEdge = layout.EdgeTop
	for i := 0; i < 20; i++ {
		childView := NewTableCell()
		childView.String = "TEST TEST"
		childView.Painter = &paint.Style{BackgroundColor: colornames.Red}
		childLayouter.Add(childView, nil)
	}

	chl10 := view.NewScrollView()
	chl10.PaintStyle = &paint.Style{BackgroundColor: colornames.Cyan}
	chl10.ContentPainter = &paint.Style{BackgroundColor: colornames.White}
	chl10.ContentLayouter = childLayouter
	chl10.ContentChildren = childLayouter.Views()
	chl10.OnScroll = func(offset layout.Point) {
		fmt.Println("scroll", offset)
	}
	_ = l.Add(chl10, func(s *constraint.Solver) {
		s.TopEqual(g4.Bottom())
		s.LeftEqual(g4.Left())
		s.Width(200)
		s.Height(200)
	})

	chl12 := view.NewSlider()
	chl12.ValueNotifier = &v.sliderValue
	chl12.MaxValue = 1
	chl12.MinValue = 0
	chl12.OnChange = func(value float64) {
		v.sliderValue.SetValue(value)
		fmt.Println("slider", value, v.sliderValue.Value())
	}
	chl12p := view.WithPainter(chl12, &paint.Style{BackgroundColor: colornames.Blue})
	g12 := l.Add(chl12p, func(s *constraint.Solver) {
		s.TopEqual(l.Top().Add(50))
		s.LeftEqual(l.Left())
		s.Width(150)
	})

	chl13 := ios.NewProgressView()
	chl13.ProgressNotifier = &v.sliderValue
	chl13.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	_ = l.Add(chl13, func(s *constraint.Solver) {
		s.TopEqual(g12.Bottom().Add(5))
		s.LeftEqual(l.Left())
		s.Width(150)
	})

	chl14 := ios.NewSegmentView()
	chl14.Value = v.segmentValue
	chl14.Titles = []string{"Title1", "Title2", "Title3"}
	chl14.OnChange = func(a int) {
		v.segmentValue = a
		fmt.Println("on value change", a)
	}
	chl14.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	_ = l.Add(chl14, func(s *constraint.Solver) {
		s.TopEqual(g5.Bottom())
		s.RightEqual(g5.Right())
		s.Width(150)
	})

	options := []view.Option{}
	if v.counter%2 == 0 {
		options = append(options, &ios.ActivityIndicator{})
		options = append(options, &ios.StatusBar{
			Hidden: true,
			Style:  ios.StatusBarStyleDark,
		})
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
		Options:  options,
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

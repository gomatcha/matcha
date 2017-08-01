package complex

import (
	"fmt"
	"strconv"
	"time"

	"golang.org/x/image/colornames"

	"gomatcha.io/bridge"
	"gomatcha.io/matcha/animate"
	"gomatcha.io/matcha/app"
	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/layout/table"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/alert"
	"gomatcha.io/matcha/view/basicview"
	"gomatcha.io/matcha/view/button"
	"gomatcha.io/matcha/view/imageview"
	"gomatcha.io/matcha/view/progressview"
	"gomatcha.io/matcha/view/scrollview"
	"gomatcha.io/matcha/view/segmentview"
	"gomatcha.io/matcha/view/slider"
	"gomatcha.io/matcha/view/switchview"
	"gomatcha.io/matcha/view/textview"
	"gomatcha.io/matcha/view/urlimageview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/complex New", func() *view.Root {
		return view.NewRoot(New(nil, ""))
	})
}

type NestedView struct {
	view.Embed
	counter      int
	sliderValue  comm.Float64Value
	segmentValue int
	value        animate.Value
}

func New(ctx *view.Context, key string) *NestedView {
	if v, ok := ctx.Prev(key).(*NestedView); ok {
		return v
	}
	return &NestedView{
		Embed: ctx.NewEmbed(key),
	}
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

func (v *NestedView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	value := animate.FloatLerp{Start: 0, End: 150}.Notifier(&v.value)
	// value := animate.FloatInterpolate(animate.FloatLerp{Start: 0, End: 150}, &v.value)

	chl1 := basicview.New(ctx, "1")
	chl1.Painter = &paint.AnimatedStyle{
		BackgroundColor: animate.RGBALerp{Start: colornames.Red, End: colornames.Yellow}.Notifier(&v.value),
	}
	g1 := l.Add(chl1, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.WidthEqual(constraint.Notifier(value))
		s.HeightEqual(constraint.Notifier(value))
	})

	chl2 := basicview.New(ctx, "2")
	chl2.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	g2 := l.Add(chl2, func(s *constraint.Solver) {
		s.TopEqual(g1.Bottom())
		s.LeftEqual(g1.Left())
		s.Width(300)
		s.Height(300)
	})

	chl3 := basicview.New(ctx, "3")
	chl3.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	g3 := l.Add(chl3, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g2.Left())
		s.Width(100)
		s.Height(100)
	})

	chl4 := basicview.New(ctx, "4")
	chl4.Painter = &paint.Style{BackgroundColor: colornames.Magenta}
	g4 := l.Add(chl4, func(s *constraint.Solver) {
		s.TopEqual(g2.Bottom())
		s.LeftEqual(g3.Right())
		s.Width(50)
		s.Height(50)
	})

	chl5 := textview.New(ctx, "a")
	chl5.String = "Subtitle"
	chl5.Style.SetAlignment(text.AlignmentCenter)
	chl5.Style.SetStrikethroughStyle(text.StrikethroughStyleSingle)
	chl5.Style.SetStrikethroughColor(colornames.Magenta)
	chl5.Style.SetUnderlineStyle(text.UnderlineStyleDouble)
	chl5.Style.SetUnderlineColor(colornames.Green)
	chl5.Style.SetFont(text.Font{
		Family: "American Typewriter",
		Face:   "Bold",
		Size:   20,
	})
	chl5p := view.WithPainter(chl5, &paint.Style{BackgroundColor: colornames.Cyan})

	g5 := l.Add(chl5p, func(s *constraint.Solver) {
		s.BottomEqual(g2.Bottom())
		s.RightEqual(g2.Right().Add(-15))
	})

	chl6 := textview.New(ctx, "6")
	chl6.String = fmt.Sprintf("Counter: %v", v.counter)
	chl6.Style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   20,
	})
	chl6p := view.WithPainter(chl6, &paint.Style{BackgroundColor: colornames.Red})
	g6 := l.Add(chl6p, func(s *constraint.Solver) {
		s.BottomEqual(g5.Top())
		s.RightEqual(g2.Right().Add(-15))
	})

	chl8 := button.New(ctx, "8")
	chl8.Text = "Button"
	chl8.OnPress = func() {
		v.counter += 1
		v.Signal()

		alert.Alert("Alert", "Message",
			&alert.Button{
				Title: "OK",
				OnPress: func() {
					fmt.Println("OK")
				},
			},
			&alert.Button{
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
		chl9 := urlimageview.New(ctx, "7")
		chl9.URL = "https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png"
		chl9.ResizeMode = imageview.ResizeModeFit
		pChl9 := view.WithPainter(chl9, &paint.Style{BackgroundColor: colornames.Cyan})

		_ = l.Add(pChl9, func(s *constraint.Solver) {
			s.BottomEqual(g8.Top())
			s.RightEqual(g2.Right().Add(-15))
			s.Width(200)
			s.Height(200)
		})
	}
	chl11 := switchview.New(ctx, "12")
	chl11.OnValueChange = func(value bool) {
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
	for i := 0; i < 20; i++ {
		childView := NewTableCell(ctx, "a"+strconv.Itoa(i))
		childView.String = "TEST TEST"
		childView.Painter = &paint.Style{BackgroundColor: colornames.Red}
		childLayouter.Add(childView, nil)
	}

	chl10 := scrollview.New(ctx, "10")
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

	chl12 := slider.New(ctx, "11")
	chl12.ValueNotifier = &v.sliderValue
	chl12.MaxValue = 1
	chl12.MinValue = 0
	chl12.OnValueChange = func(value float64) {
		v.sliderValue.SetValue(value)
		fmt.Println("slider", value, v.sliderValue.Value())
	}
	chl12p := view.WithPainter(chl12, &paint.Style{BackgroundColor: colornames.Blue})
	g12 := l.Add(chl12p, func(s *constraint.Solver) {
		s.TopEqual(l.Top().Add(50))
		s.LeftEqual(l.Left())
		s.Width(150)
	})

	chl13 := progressview.New(ctx, "13")
	chl13.ProgressNotifier = &v.sliderValue
	chl13.PaintStyle = &paint.Style{BackgroundColor: colornames.White}
	_ = l.Add(chl13, func(s *constraint.Solver) {
		s.TopEqual(g12.Bottom().Add(5))
		s.LeftEqual(l.Left())
		s.Width(150)
	})

	chl14 := segmentview.New(ctx, "14")
	chl14.Value = v.segmentValue
	chl14.Titles = []string{"Title1", "Title2", "Title3"}
	chl14.OnValueChange = func(a int) {
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
		options = append(options, app.ActivityIndicator{})
		options = append(options, app.StatusBar{
			Hidden: true,
			Style:  app.StatusBarStyleDark,
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

func NewTableCell(ctx *view.Context, key string) *TableCell {
	if v, ok := ctx.Prev(key).(*TableCell); ok {
		return v
	}
	return &TableCell{
		Embed: ctx.NewEmbed(key),
	}
}

func (v *TableCell) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	l.Solve(func(s *constraint.Solver) {
		s.Height(50)
	})

	textView := textview.New(ctx, "1")
	textView.String = v.String
	textView.Style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   20,
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

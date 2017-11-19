package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/view NewTextView", func() view.View {
		return NewTextView()
	})
}

type TextViewTest struct {
	view.Embed
	text      *text.Text
	responder *keyboard.Responder
}

func NewTextView() *TextViewTest {
	return &TextViewTest{
		text:      text.New("blah"),
		responder: &keyboard.Responder{},
	}
}

func (v *TextViewTest) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageVisible) {
		v.responder.Show()
	}
}

func (v *TextViewTest) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	label := view.NewTextView()
	label.String = "Text view:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(15)
		s.Left(15)
	})

	textView := view.NewTextView()
	textView.String = "Test1234"
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Single Line:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.MaxLines = 1
	textView.String = "Line 1\nLine 2\nLine 3\n Line 4"
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Multi Line:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.MaxLines = 3
	textView.String = "Line 1\nLine 2\nLine 3\n Line 4"
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Font:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetFont(text.DefaultBoldFont(24))
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.RightEqual(l.Right().Add(-15))
	})

	label = view.NewTextView()
	label.String = "Left align:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetAlignment(text.AlignmentLeft)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.RightEqual(l.Right().Add(-15))
	})

	label = view.NewTextView()
	label.String = "Center align:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetAlignment(text.AlignmentCenter)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.RightEqual(l.Right().Add(-15))
	})

	label = view.NewTextView()
	label.String = "Right align:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetAlignment(text.AlignmentRight)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.RightEqual(l.Right().Add(-15))
	})

	label = view.NewTextView()
	label.String = "Text color:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetTextColor(colornames.Red)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Strikethrough:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetStrikethroughStyle(text.StrikethroughStyleSingle)
	textView.Style.SetStrikethroughColor(colornames.Red)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Underline:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.String = "Test1234"
	textView.Style.SetUnderlineStyle(text.UnderlineStyleSingle)
	textView.Style.SetUnderlineColor(colornames.Red)
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Text:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	textView = view.NewTextView()
	textView.Text = text.New("Test1234")
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	label = view.NewTextView()
	label.String = "Styled Text:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	style := &text.Style{}
	style.SetTextColor(colornames.Red)
	style.SetFont(text.DefaultBoldFont(24))

	style2 := &text.Style{}
	style2.SetTextColor(colornames.Green)
	style2.SetFont(text.DefaultItalicFont(18))

	st := text.NewStyledText("Test1234", style)
	st.Set(style2, 4, 8)

	textView = view.NewTextView()
	textView.StyledText = st
	g = l.Add(textView, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	/*
		style := &text.Style{}
		style.SetAlignment(text.AlignmentLeft)
		style.SetStrikethroughStyle(text.StrikethroughStyleDouble)
		style.SetStrikethroughColor(colornames.Blue)
		style.SetUnderlineStyle(text.UnderlineStyleDouble)
		style.SetUnderlineColor(colornames.Green)
		style.SetTextColor(colornames.Yellow)
		style.SetFont(text.FontWithName("HelveticaNeue-Bold", 20))

		style2 := &text.Style{}
		style2.SetAlignment(text.AlignmentLeft)
		style2.SetStrikethroughStyle(text.StrikethroughStyleDouble)
		style2.SetStrikethroughColor(colornames.Blue)
		style2.SetUnderlineStyle(text.UnderlineStyleDouble)
		style2.SetUnderlineColor(colornames.Green)
		style2.SetTextColor(colornames.Red)
		style2.SetFont(text.FontWithName("HelveticaNeue", 10))

		st := text.NewStyledText("Subtitle", style)
		st.Set(style2, 0, 3)

		chl := view.NewTextView()
		chl.StyledText = st
		chlP := view.WithPainter(chl, &paint.Style{BackgroundColor: colornames.Blue})
		chlG := l.Add(chlP, func(s *constraint.Solver) {
			s.TopEqual(constraint.Const(50))
			s.LeftEqual(constraint.Const(100))
			s.Width(200)
		})

		reverse := view.NewTextView()
		reverse.Style.SetAlignment(text.AlignmentCenter)
		reverse.Style.SetStrikethroughStyle(text.StrikethroughStyleDouble)
		reverse.Style.SetStrikethroughColor(colornames.Blue)
		reverse.Style.SetUnderlineStyle(text.UnderlineStyleDouble)
		reverse.Style.SetUnderlineColor(colornames.Green)
		reverse.Style.SetTextColor(colornames.Black)
		reverse.Style.SetFont(text.FontWithName("HelveticaNeue-Bold", 20))
		reverse.String = Reverse(v.text.String())
		reverse.PaintStyle = &paint.Style{BackgroundColor: colornames.Green}
		l.Add(reverse, func(s *constraint.Solver) {
			s.TopEqual(chlG.Bottom())
			s.LeftEqual(chlG.Left())
		})

		button1 := view.NewButton()
		// button1.PaintStyle = &paint.Style{BackgroundColor: colornames.Blue}
		// button1.Color = colornames.Green
		button1.String = "Toggle Keyboard"
		button1.OnPress = func() {
			if !v.responder.Visible() {
				v.responder.Show()
			} else {
				v.responder.Dismiss()
			}
		}
		l.Add(button1, func(s *constraint.Solver) {
			s.Top(200)
			s.Left(100)
		})

		input := view.NewTextInput()
		input.RWText = v.text
		input.Placeholder = "Placeholder"
		input.KeyboardType = keyboard.URLType
		input.MaxLines = 1
		input.Responder = v.responder
		input.OnChange = func(t *text.Text) {
			v.Signal()
		}
		input.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
		l.Add(input, func(s *constraint.Solver) {
			s.Top(100)
			s.Left(100)
			s.Width(200)
			s.Height(100)
		})
	*/
	return view.Model{
		Children: l.Views(),
		Painter:  &paint.Style{BackgroundColor: colornames.White},
		Layouter: l,
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

package view

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
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

type TextView struct {
	view.Embed
	text      *text.Text
	responder *keyboard.Responder
}

func NewTextView() *TextView {
	return &TextView{
		text:      text.New("blah"),
		responder: &keyboard.Responder{},
	}
}

func (v *TextView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageVisible) {
		v.responder.Show()
	}
}

func (v *TextView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

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
	input.Text = v.text
	input.Placeholder = "Placeholder"
	input.KeyboardType = keyboard.URLType
	input.MaxLines = 1
	input.KeyboardReturnType = keyboard.SearchReturnType
	input.Responder = v.responder
	input.OnTextChange = func(t *text.Text) {
		v.Signal()
	}
	input.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
	l.Add(input, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(100)
		s.Width(200)
		s.Height(100)
	})

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

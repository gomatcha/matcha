// Package textview provides examples of how to use the matcha/view/textview package.
package textview

import (
	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/button"
	"gomatcha.io/matcha/view/textinput"
	"gomatcha.io/matcha/view/textview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/textview New", func() *view.Root {
		return view.NewRoot(New())
	})
}

type TextView struct {
	view.Embed
	text      *text.Text
	responder *keyboard.Responder
}

func New() *TextView {
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

func (v *TextView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	style := &text.Style{}
	style.SetAlignment(text.AlignmentCenter)
	style.SetStrikethroughStyle(text.StrikethroughStyleDouble)
	style.SetStrikethroughColor(colornames.Blue)
	style.SetUnderlineStyle(text.UnderlineStyleDouble)
	style.SetUnderlineColor(colornames.Green)
	style.SetTextColor(colornames.Yellow)
	style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Face:   "Bold",
		Size:   20,
	})

	style2 := &text.Style{}
	style2.SetAlignment(text.AlignmentCenter)
	style2.SetStrikethroughStyle(text.StrikethroughStyleDouble)
	style2.SetStrikethroughColor(colornames.Blue)
	style2.SetUnderlineStyle(text.UnderlineStyleDouble)
	style2.SetUnderlineColor(colornames.Green)
	style2.SetTextColor(colornames.Red)
	style2.SetFont(text.Font{
		Family: "Helvetica Neue",
		Size:   10,
	})

	st := text.NewStyledText("Subtitle", style)
	st.Set(style2, 0, 3)

	chl := textview.New()
	chl.StyledText = st
	chlP := view.WithPainter(chl, &paint.Style{BackgroundColor: colornames.Blue})
	chlG := l.Add(chlP, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(100))
		s.LeftEqual(constraint.Const(100))
	})

	reverse := textview.New()
	reverse.Style.SetAlignment(text.AlignmentCenter)
	reverse.Style.SetStrikethroughStyle(text.StrikethroughStyleDouble)
	reverse.Style.SetStrikethroughColor(colornames.Blue)
	reverse.Style.SetUnderlineStyle(text.UnderlineStyleDouble)
	reverse.Style.SetUnderlineColor(colornames.Green)
	reverse.Style.SetTextColor(colornames.Black)
	reverse.Style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Face:   "Bold",
		Size:   20,
	})
	reverse.String = Reverse(v.text.String())
	l.Add(reverse, func(s *constraint.Solver) {
		s.TopEqual(chlG.Bottom())
		s.LeftEqual(chlG.Left())
	})

	button1 := button.New()
	button1.Text = "Toggle Keyboard"
	button1.OnPress = func() {
		if !v.responder.Visible() {
			v.responder.Show()
		} else {
			v.responder.Dismiss()
		}
	}
	l.Add(button1, func(s *constraint.Solver) {
		s.Top(300)
		s.Left(100)
	})

	input := textinput.New()
	input.Text = v.text
	input.KeyboardType = keyboard.URLType
	input.KeyboardAppearance = keyboard.DarkAppearance
	input.KeyboardReturnType = keyboard.GoogleReturnType
	input.Responder = v.responder
	input.OnTextChange = func(t *text.Text) {
		v.Signal()
	}
	input.PaintStyle = &paint.Style{BackgroundColor: colornames.Lightgray}
	l.Add(input, func(s *constraint.Solver) {
		s.Top(200)
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

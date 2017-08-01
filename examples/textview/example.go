package textview

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/bridge"
	"gomatcha.io/matcha/keyboard"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/touch"
	"gomatcha.io/matcha/view"
	"gomatcha.io/matcha/view/textinput"
	"gomatcha.io/matcha/view/textview"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/textview New", func() *view.Root {
		return view.NewRoot(New(nil, ""))
	})
}

type TextView struct {
	view.Embed
	text      *text.Text
	responder *keyboard.Responder
}

func New(ctx *view.Context, key string) *TextView {
	if v, ok := ctx.Prev(key).(*TextView); ok {
		return v
	}
	return &TextView{
		Embed:     ctx.NewEmbed(key),
		text:      text.New("blah"),
		responder: &keyboard.Responder{},
	}
}

func (v *TextView) Lifecycle(from, to view.Stage) {
	if view.EntersStage(from, to, view.StageVisible) {
		v.responder.Show()
		fmt.Println("show", v.responder.Visible())
	}
}

func (v *TextView) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}

	chl := textview.New(ctx, "a")
	chl.String = "Subtitle"
	chl.Style.SetAlignment(text.AlignmentCenter)
	chl.Style.SetStrikethroughStyle(text.StrikethroughStyleDouble)
	chl.Style.SetStrikethroughColor(colornames.Blue)
	chl.Style.SetUnderlineStyle(text.UnderlineStyleDouble)
	chl.Style.SetUnderlineColor(colornames.Blue)
	chl.Style.SetTextColor(colornames.Yellow)
	chl.Style.SetFont(text.Font{
		Family: "Helvetica Neue",
		Face:   "Bold",
		Size:   20,
	})
	chlP := view.WithPainter(chl, &paint.Style{BackgroundColor: colornames.Blue})
	chlG := l.Add(chlP, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(100))
		s.LeftEqual(constraint.Const(100))
	})

	input := textinput.New(ctx, "input")
	input.Text = v.text
	input.KeyboardType = keyboard.URLType
	input.KeyboardAppearance = keyboard.DarkAppearance
	input.KeyboardReturnType = keyboard.GoogleReturnType
	input.Responder = v.responder
	input.OnTextChange = func(t *text.Text) {
		v.Signal()
	}
	inputP := view.WithPainter(input, &paint.Style{BackgroundColor: colornames.Yellow})
	l.Add(inputP, func(s *constraint.Solver) {
		s.TopEqual(constraint.Const(200))
		s.LeftEqual(constraint.Const(100))
		s.WidthEqual(constraint.Const(200))
		s.HeightEqual(constraint.Const(100))
	})

	reverse := textview.New(ctx, "reverse")
	reverse.String = Reverse(v.text.String())
	l.Add(reverse, func(s *constraint.Solver) {
		s.TopEqual(chlG.Bottom())
		s.LeftEqual(chlG.Left())
	})

	tap := &touch.TapRecognizer{
		Count: 2,
		OnTouch: func(e *touch.TapEvent) {
			v.responder.Dismiss()
		},
	}

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
		Options: []view.Option{
			touch.RecognizerList{tap},
		},
	}
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Package pointer provides examples of how to use the matcha/pointer package.
package pointer

import (
	"fmt"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pointer"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/pointer NewTouchView", func() view.View {
		return NewTouchView()
	})
}

type TouchView struct {
	view.Embed
	tapCounter              int
	tap2Counter             int
	tap11Counter            int
	tap22Counter            int
	pressCounter            int
	buttonCounter           int
	buttonChildCounter      int
	scrollButtonCounter     int
	scrollButtonHighlighted bool
}

func NewTouchView() *TouchView {
	return &TouchView{
		Embed: view.Embed{},
	}
}

func (v *TouchView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	tap := NewTapChildView()
	tap.Count = 1
	tap.OnTouch = func() {
		v.tapCounter += 1
		go v.Signal() // TODO(KD): Why is this on separate thread?
	}
	g := l.Add(tap, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(100)
		s.Height(100)
	})

	tapCount := view.NewTextView()
	tapCount.String = fmt.Sprintf("tap: %v", v.tapCounter)
	tapCount.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(tapCount, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	tap2 := NewTapChildView()
	tap2.Count = 2
	tap2.OnTouch = func() {
		v.tap2Counter += 1
		go v.Signal() // TODO(KD): Why is this on separate thread?
	}
	g = l.Add(tap2, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	tap2Count := view.NewTextView()
	tap2Count.String = fmt.Sprintf("double tap: %v", v.tap2Counter)
	tap2Count.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(tap2Count, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	tap12 := view.NewBasicView()
	tap12.Painter = &paint.Style{BackgroundColor: colornames.Blue}
	tap12.Options = []view.Option{
		&pointer.TapGesture{
			Count: 1,
			OnRecognize: func(e *pointer.TapEvent) {
				v.tap11Counter += 1
				v.Signal()
			},
		},
		&pointer.TapGesture{
			Count: 2,
			OnRecognize: func(e *pointer.TapEvent) {
				v.tap22Counter += 1
				v.Signal()
			},
		},
	}
	g = l.Add(tap12, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	count := view.NewTextView()
	count.String = fmt.Sprintf("single tap: %v", v.tap11Counter)
	count.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(count, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	count = view.NewTextView()
	count.String = fmt.Sprintf("double tap: %v", v.tap22Counter)
	count.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(count, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	button := NewButtonChildView()
	button.OnTouch = func() {
		fmt.Println("On Button")
		v.buttonCounter += 1
		go v.Signal()
	}
	button.OnChildTouch = func() {
		fmt.Println("On Child Button")
		v.buttonChildCounter += 1
		go v.Signal()
	}
	g = l.Add(button, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Width(100)
		s.Height(100)
	})

	buttonOverlay := view.NewBasicView()
	buttonOverlay.Painter = &paint.Style{BackgroundColor: colornames.Yellow}
	l.Add(buttonOverlay, func(s *constraint.Solver) {
		s.TopEqual(g.Top())
		s.LeftEqual(g.Left())
		s.Width(50)
		s.Height(50)
	})

	buttonCount := view.NewTextView()
	buttonCount.String = fmt.Sprintf("button: %v", v.buttonCounter)
	buttonCount.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(buttonCount, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	buttonChild := view.NewTextView()
	buttonChild.String = fmt.Sprintf("button child: %v", v.buttonChildCounter)
	buttonChild.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(buttonChild, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	scrollButtonColor := colornames.Blue
	if v.scrollButtonHighlighted {
		scrollButtonColor = colornames.Red
	}
	scrollButton := view.NewBasicView()
	scrollButton.Painter = &paint.Style{BackgroundColor: scrollButtonColor}
	scrollButton.Options = []view.Option{
		&pointer.ButtonGesture{
			OnHighlight: func(highlighted bool) {
				v.scrollButtonHighlighted = highlighted
				v.Signal()
			},
			OnRecognize: func(e *pointer.ButtonEvent) {
				v.scrollButtonHighlighted = false
				v.scrollButtonCounter += 1
			},
		},
	}

	scrollLayouter := &constraint.Layouter{}
	scrollLayouter.Solve(func(s *constraint.Solver) {
		s.Width(100)
		s.Height(500)
	})
	scrollLayouter.Add(scrollButton, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(0)
		s.Width(100)
		s.Height(50)
	})

	scrollview := view.NewScrollView()
	scrollview.ContentChildren = scrollLayouter.Views()
	scrollview.ContentPainter = &paint.Style{BackgroundColor: colornames.Yellow}
	scrollview.ContentLayouter = scrollLayouter
	g = l.Add(scrollview, func(s *constraint.Solver) {
		s.Top(100)
		s.Left(200)
		s.Width(100)
		s.Height(100)
	})

	scrollButtonCount := view.NewTextView()
	scrollButtonCount.String = fmt.Sprintf("button: %v", v.scrollButtonCounter)
	scrollButtonCount.Style.SetFont(text.FontWithName("HelveticaNeue", 20))
	g = l.Add(scrollButtonCount, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Green},
	}
}

type TapChildView struct {
	view.Embed
	OnTouch func()
	Count   int
}

func NewTapChildView() *TapChildView {
	return &TapChildView{}
}

func (v *TapChildView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}
	// a := view.NewBasicView()
	// a.Painter = &paint.Style{BackgroundColor: colornames.Red}
	// l.Add(a, func(s *constraint.Solver) {
	// 	s.Left(0)
	// 	s.Top(0)
	// 	s.Width(100)
	// 	s.Height(100)
	// })
	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.Blue},
		Options: []view.Option{
			&pointer.TapGesture{
				Count: v.Count,
				OnRecognize: func(e *pointer.TapEvent) {
					fmt.Println("Tap Recognized")
					v.OnTouch()
				},
			},
		},
	}
}

// type PressChildView struct {
// 	view.Embed
// 	OnTouch func()
// }

// func NewPressChildView() *PressChildView {
// 	return &PressChildView{}
// }

// func (v *PressChildView) Build(ctx view.Context) view.Model {
// 	return view.Model{
// 		Painter: &paint.Style{BackgroundColor: colornames.Blue},
// 		Options: []view.Option{
// 			pointer.GestureList{&pointer.PressGesture{
// 				MinDuration: time.Second / 2,
// 				OnEvent: func(e *pointer.PressEvent) {
// 					if e.Kind == pointer.EventKindPossible {
// 						fmt.Println("Press Possible")
// 					} else if e.Kind == pointer.EventKindChanged {
// 						fmt.Println("Press Changed")
// 					} else if e.Kind == pointer.EventKindFailed {
// 						fmt.Println("Press Failed")
// 					} else if e.Kind == pointer.EventKindRecognized {
// 						fmt.Println("Press Recognized")
// 						v.OnTouch()
// 					}
// 				},
// 			}},
// 		},
// 	}
// }

type ButtonChildView struct {
	view.Embed
	OnTouch          func()
	OnChildTouch     func()
	highlighted      bool
	childHighlighted bool
}

func NewButtonChildView() *ButtonChildView {
	return &ButtonChildView{}
}

func (v *ButtonChildView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	color := colornames.Blue
	if v.highlighted {
		color = colornames.Red
	}

	childColor := colornames.Purple
	if v.childHighlighted {
		childColor = colornames.Pink
	}

	child := view.NewBasicView()
	child.Painter = &paint.Style{BackgroundColor: childColor}
	child.Options = []view.Option{
		&pointer.ButtonGesture{
			// Exclusive: true,
			OnHighlight: func(highlighted bool) {
				v.childHighlighted = highlighted
				v.Signal()
			},
			OnRecognize: func(e *pointer.ButtonEvent) {
				v.childHighlighted = false
				v.OnChildTouch()
			},
		},
	}
	l.Add(child, func(s *constraint.Solver) {
		s.Top(0)
		s.Left(50)
		s.Width(50)
		s.Height(50)
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: color},
		Options: []view.Option{
			&pointer.ButtonGesture{
				OnHighlight: func(highlighted bool) {
					v.highlighted = highlighted
					v.Signal()
				},
				OnRecognize: func(e *pointer.ButtonEvent) {
					v.highlighted = false
					v.OnTouch()
				},
			},
		},
	}
}

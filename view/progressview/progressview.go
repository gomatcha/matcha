// Package progressview implements a native progress view.
package progressview

import (
	"image/color"

	"gomatcha.io/matcha/comm"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/pb"
	"gomatcha.io/matcha/pb/view/progressview"
	"gomatcha.io/matcha/view"
)

// View implements a progess view.
type View struct {
	view.Embed
	Progress         float64
	ProgressNotifier comm.Float64Notifier
	ProgressColor    color.Color
	PaintStyle       *paint.Style
	progressNotifier comm.Float64Notifier
}

// New returns either the previous View in ctx with matching key, or a new View if none exists.
func New() *View {
	return &View{}
}

// Lifecycle implements the view.View interface.
func (v *View) Lifecycle(from, to view.Stage) {
	if view.ExitsStage(from, to, view.StageMounted) {
		if v.progressNotifier != nil {
			v.Unsubscribe(v.progressNotifier)
		}
	}
}

// Build implements the view.View interface.
func (v *View) Build(ctx *view.Context) view.Model {
	l := &constraint.Layouter{}
	l.Solve(func(s *constraint.Solver) {
		s.Height(2) // 2.5 if its a bar
		s.WidthEqual(l.MinGuide().Width())
		s.TopEqual(l.MaxGuide().Top())
		s.LeftEqual(l.MaxGuide().Left())
	})

	if v.ProgressNotifier != v.progressNotifier {
		if v.progressNotifier != nil {
			v.Unsubscribe(v.progressNotifier)
		}
		if v.ProgressNotifier != nil {
			v.Subscribe(v.ProgressNotifier)
		}
		v.progressNotifier = v.ProgressNotifier
	}

	val := v.Progress
	if v.ProgressNotifier != nil {
		val = v.ProgressNotifier.Value()
	}

	painter := paint.Painter(nil)
	if v.PaintStyle != nil {
		painter = v.PaintStyle
	}
	return view.Model{
		Painter:        painter,
		Layouter:       l,
		NativeViewName: "gomatcha.io/matcha/view/progressview",
		NativeViewState: &progressview.View{
			Progress:      val,
			ProgressColor: pb.ColorEncode(v.ProgressColor),
		},
	}
}

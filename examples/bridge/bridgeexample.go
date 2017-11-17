package bridge

import (
	"fmt"
	"runtime"

	"golang.org/x/image/colornames"
	"gomatcha.io/matcha/bridge"
	"gomatcha.io/matcha/layout/constraint"
	"gomatcha.io/matcha/paint"
	"gomatcha.io/matcha/text"
	"gomatcha.io/matcha/view"
)

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/bridge NewBridgeView", func() view.View {
		return NewBridgeView()
	})
}

func init() {
	bridge.RegisterFunc("gomatcha.io/matcha/examples/bridge callWithGoValues", func(v string) string {
		return fmt.Sprintf("Goodbye %v!", v)
	})
	bridge.RegisterFunc("gomatcha.io/matcha/examples/bridge callWithForeignValues", func(v *bridge.Value) *bridge.Value {
		return bridge.String(fmt.Sprintf("Goodbye %v!", v.ToString()))
	})
}

type BridgeView struct {
	view.Embed
}

func NewBridgeView() *BridgeView {
	return &BridgeView{}
}

func (v *BridgeView) Build(ctx view.Context) view.Model {
	l := &constraint.Layouter{}

	// Get the corresponding native bridge object. See ExampleObjcBridge.m and ExampleJavaBridge.java.
	brg := bridge.Bridge("gomatcha.io/matcha/example/bridge")

	var str string
	// Call the method with a Go object as a parameter.
	if runtime.GOOS == "android" { // Android and iOS have different method signatures.
		str = brg.Call("callWithGoValues", bridge.Interface("Ame")).ToInterface().(string)
	} else {
		str = brg.Call("callWithGoValues:", bridge.Interface("Ame")).ToInterface().(string)
	}

	label := view.NewTextView()
	label.String = "Calling a objc/java method from Go, passing Go values:"
	label.Style.SetFont(text.DefaultFont(18))
	g := l.Add(label, func(s *constraint.Solver) {
		s.Top(50)
		s.Left(15)
		s.Right(-15)
	})

	label = view.NewTextView()
	label.String = str
	label.Style.SetFont(text.DefaultFont(15))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	// Call the method with a foreign(objc/java) object as a parameter.
	if runtime.GOOS == "android" {
		str = brg.Call("callWithForeignValues", bridge.String("Yuki")).ToString()
	} else {
		str = brg.Call("callWithForeignValues:", bridge.String("Yuki")).ToString()
	}

	label = view.NewTextView()
	label.String = "Calling a objc/java method from Go, passing objc/java values:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Right(-15)
	})

	label = view.NewTextView()
	label.String = str
	label.Style.SetFont(text.DefaultFont(15))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	// Call foreign function, which in turn calls the `gomatcha.io/matcha/examples/bridge callWithForeignValues` function we registered in init().
	if runtime.GOOS == "android" {
		str = brg.Call("callGoFunctionWithForeignValues").ToString()
	} else {
		str = brg.Call("callGoFunctionWithForeignValues").ToString()
	}

	label = view.NewTextView()
	label.String = "Calling a Go function from objc/java code, passing objc/java values:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Right(-15)
	})

	label = view.NewTextView()
	label.String = str
	label.Style.SetFont(text.DefaultFont(15))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	// Call foreign function, which in turn calls the `gomatcha.io/matcha/examples/bridge callWithGoValues` function we registered in init().
	if runtime.GOOS == "android" {
		str = brg.Call("callGoFunctionWithGoValues").ToString()
	} else {
		str = brg.Call("callGoFunctionWithGoValues").ToString()
	}

	label = view.NewTextView()
	label.String = "Calling a Go function from objc/java code, passing Go values:"
	label.Style.SetFont(text.DefaultFont(18))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
		s.Right(-15)
	})

	label = view.NewTextView()
	label.String = str
	label.Style.SetFont(text.DefaultFont(15))
	g = l.Add(label, func(s *constraint.Solver) {
		s.TopEqual(g.Bottom())
		s.LeftEqual(g.Left())
	})

	return view.Model{
		Children: l.Views(),
		Layouter: l,
		Painter:  &paint.Style{BackgroundColor: colornames.White},
	}
}

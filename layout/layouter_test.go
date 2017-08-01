package layout

import (
	"testing"
)

func TestFit(t *testing.T) {
	g := Guide{Frame: Rt(0, 0, 100, 100)}
	ctx := &Context{MinSize: Pt(50, 50), MaxSize: Pt(75, 75)}

	fitG := g.Fit(ctx)
	expect := Guide{Frame: Rt(0, 0, 75, 75)}
	if fitG != expect {
		t.Error("Error")
	}
}

func TestFit2(t *testing.T) {
	g := Guide{Frame: Rt(0, 0, 25, 25)}
	ctx := &Context{MinSize: Pt(50, 50), MaxSize: Pt(75, 75)}

	fitG := g.Fit(ctx)
	expect := Guide{Frame: Rt(0, 0, 50, 50)}
	if fitG != expect {
		t.Error("Error")
	}
}

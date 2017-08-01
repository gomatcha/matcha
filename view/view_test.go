package view

import "testing"

func TestStage(t *testing.T) {
	test := []struct {
		from, to, test int
		enter, exit    bool
	}{
		{1, 1, 3, false, false},
		{1, 2, 3, false, false},
		{1, 3, 3, true, false},
		{1, 4, 3, true, false},
		{1, 5, 3, true, false},

		{2, 1, 3, false, false},
		{2, 2, 3, false, false},
		{2, 3, 3, true, false},
		{2, 4, 3, true, false},
		{2, 5, 3, true, false},

		{3, 1, 3, false, true},
		{3, 2, 3, false, true},
		{3, 3, 3, false, false},
		{3, 4, 3, false, false},
		{3, 5, 3, false, false},

		{4, 1, 3, false, true},
		{4, 2, 3, false, true},
		{4, 3, 3, false, false},
		{4, 4, 3, false, false},
		{4, 5, 3, false, false},

		{5, 1, 3, false, true},
		{5, 2, 3, false, true},
		{5, 3, 3, false, false},
		{5, 4, 3, false, false},
		{5, 5, 3, false, false},
	}

	for _, i := range test {
		if EntersStage(Stage(i.from), Stage(i.to), Stage(i.test)) != i.enter {
			t.Error("enter", i.from, i.to, i.test)
		}
		if ExitsStage(Stage(i.from), Stage(i.to), Stage(i.test)) != i.exit {
			t.Error("exit", i.from, i.to, i.test)
		}
	}
}

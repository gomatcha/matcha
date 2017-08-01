package constraint

import (
	"math"
	"testing"
)

func TestConstrainedRect(t *testing.T) {
	cr := newConstrainedRect()
	if !cr.isValid() {
		t.Error("NewConstrainedRect is invalid")
	}
}

func TestIntersect(t *testing.T) {
	r := _range{0, 10}
	if n := r.intersectMin(math.Inf(-1)); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMin(-5); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMin(0); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMin(5); n != (_range{5, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMin(15); n != (_range{15, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}

	if n := r.intersectMax(-5); n != (_range{0, -5}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMax(5); n != (_range{0, 5}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMax(10); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMax(15); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersectMax(math.Inf(1)); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}

	if n := r.intersect(_range{-5, -5}); n != (_range{0, -5}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersect(_range{-5, 0}); n != (_range{0, 0}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersect(_range{-5, 5}); n != (_range{0, 5}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersect(_range{-5, 10}); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersect(_range{-5, 15}); n != (_range{0, 10}) {
		t.Errorf("Incorrect result: %v", n)
	}
	if n := r.intersect(_range{0, -5}); n != (_range{0, -5}) {
		t.Errorf("Incorrect result: %v", n)
	}
}

func TestIsValid(t *testing.T) {
	if b := (_range{0, 10}).isValid(); !b {
		t.Errorf("Incorrect result: %v", b)
	}
	if b := (_range{0, 0}).isValid(); !b {
		t.Errorf("Incorrect result: %v", b)
	}
	if b := (_range{math.Inf(1), math.Inf(1)}).isValid(); !b {
		t.Errorf("Incorrect result: %v", b)
	}
	if b := (_range{math.Inf(-1), math.Inf(1)}).isValid(); !b {
		t.Errorf("Incorrect result: %v", b)
	}
	if b := (_range{math.Inf(1), math.Inf(-1)}).isValid(); b {
		t.Errorf("Incorrect result: %v", b)
	}
	if b := (_range{10, 0}).isValid(); b {
		t.Errorf("Incorrect result: %v", b)
	}
}

func TestNearest(t *testing.T) {
	r := _range{0, 10}
	if n := r.nearest(100); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(10); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(math.Inf(1)); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(8); n != 8 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(-10); n != 0 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(math.Inf(-1)); n != 0 {
		t.Errorf("Incorrect nearest: %v", n)
	}

	// Reversed range.
	r = _range{10, 0}
	if n := r.nearest(100); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(10); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(math.Inf(1)); n != 10 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(8); n != 8 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(-10); n != 0 {
		t.Errorf("Incorrect nearest: %v", n)
	}
	if n := r.nearest(math.Inf(-1)); n != 0 {
		t.Errorf("Incorrect nearest: %v", n)
	}
}

func TestSolveWidth(t *testing.T) {
	cr := newConstrainedRect()
	if w, ok := cr.solveWidth(10); w != 10 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	if w, ok := cr.solveWidth(-10); w != 0 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	if w, ok := cr.solveWidth(5); w != 5 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	cr.centerX = _range{0, 10}
	if w, ok := cr.solveWidth(5); w != 5 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	cr.centerX = _range{3, 3}
	cr.left = _range{0, 0}
	if w, ok := cr.solveWidth(6); w != 6 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	cr.right = _range{4, 10}
	cr.left = _range{-5, 0}
	if w, ok := cr.solveWidth(2); w != 4 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	cr.right = _range{-5, 0}
	cr.left = _range{-5, 0}
	if w, ok := cr.solveWidth(15); w != 5 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}

	cr = newConstrainedRect()
	cr.width = _range{0, 10}
	cr.right = _range{-5, 0}
	cr.left = _range{-5, 0}
	cr.centerX = _range{-4, -3}
	if w, ok := cr.solveWidth(15); w != 4 || !ok.isValid() {
		t.Errorf("Incorrect solution: (%v, %v)", w, ok)
	}
}

package lang

import (
	"testing"
)

func TestIf(t *testing.T) {
	if If(true, "yes", "no") != "yes" {
		t.Error("expected yes, got no")
	}

	if If(false, "yes", "no") != "no" {
		t.Error("expected no, got yes")
	}

	yesCalled := false
	noCalled := false

	yesFn := func() string {
		yesCalled = true
		return "yes"
	}

	noFn := func() string {
		noCalled = true
		return "no"
	}

	if IfLazy(true, yesFn, noFn) != "yes" {
		t.Error("expected yes, got no")
	}

	if !yesCalled {
		t.Error("impossible to reach")
	}

	if noCalled {
		t.Error("not lazy")
	}

	yesCalled = false
	noCalled = false

	if IfLazy(false, yesFn, noFn) != "no" {
		t.Error("expected no, got yes")
	}

	if yesCalled {
		t.Error("not lazy")
	}

	if !noCalled {
		t.Error("impossible to reach")
	}
}

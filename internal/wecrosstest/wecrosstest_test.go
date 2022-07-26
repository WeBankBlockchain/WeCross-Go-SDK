package wecrosstest

import (
	"reflect"
	"testing"
)

type tRunST struct {
	setup, test, teardown bool
}

func (t *tRunST) Setup(*testing.T) {
	t.setup = true
}
func (t *tRunST) TestSubTest(*testing.T) {
	t.test = true
}
func (t *tRunST) Teardown(*testing.T) {
	t.teardown = true
}

func TestRunSubTests(t *testing.T) {
	x := &tRunST{}
	RunSubTests(t, x)
	if want := (&tRunST{setup: true, test: true, teardown: true}); !reflect.DeepEqual(x, want) {
		t.Fatalf("x = %v; want all fields true", x)
	}
}

type tNoST struct {
	test bool
}

func (t *tNoST) TestSubTest(*testing.T) {
	t.test = true
}

func TestNoSetupOrTeardown(t *testing.T) {
	// Ensures nothing panics or fails if Setup/Teardown are omitted.
	x := &tNoST{}
	RunSubTests(t, x)
	if want := (&tNoST{test: true}); !reflect.DeepEqual(x, want) {
		t.Fatalf("x = %v; want %v", x, want)
	}
}

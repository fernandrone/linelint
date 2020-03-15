package linter

import (
	"testing"
)

func TestShouldIgnore_True(t *testing.T) {
	f := ".git/objects/04/9f2973ffc85f71da1fd5a"
	got := NewEndOfFileRule().ShouldIgnore(f)

	want := true

	if got != want {
		t.Errorf("NewEndOfFileRule().ShouldIgnore(%q):\n\tExpected %v, got %v", f, want, got)
	}
}

func TestShouldIgnore_False(t *testing.T) {
	f := "README"
	got := NewEndOfFileRule().ShouldIgnore(f)

	want := false

	if got != want {
		t.Errorf("NewEndOfFileRule().ShouldIgnore(%q):\n\tExpected %v, got %v", f, want, got)
	}
}

package sexp

import "testing"

func TestCarCdr(t *testing.T) {
	var s Sexp

	// List
	s = List{Symbol("a"), Symbol("b"), Symbol("c")}
	car := s.Head()
	if car != Symbol("a") {
		t.Error("Expected .Head() to be \"a\"")
	}

	cdr := s.Tail()
	l := cdr.(List)
	if len(l) != 2 {
		t.Fatal("Expected .Tail to have 2 elem")
	}
	if l[0] != Symbol("b") {
		t.Error("Expected cdr[0] to be \"b\"")
	}

	if l[1] != Symbol("c") {
		t.Error("Expected cdr[1] to be \"c\"")
	}

	// Atom
	s = Symbol("a")
	if s.Head() != Symbol("a") {
		t.Errorf("Expected .Head() to return itself for atoms")
	}

	if s.Tail() != nil {
		t.Error("Atoms should have nil .Tail()s")
	}

	// Strict
	s = NewStrict(Symbol("a"))
	s = addChild(s, Symbol("b"), true)
	s = addChild(s, Symbol("c"), true)

	if s.Head() != Symbol("a") {
		t.Error("Expected .Head() to be \"a\"")
	}

	cdr = s.Tail()
	if cdr.Head() != Symbol("b") {
		t.Errorf("Expected cdr[0] to be \"b\"")
	}

	if cdr.Tail().Head() != Symbol("c") {
		t.Error("Expected cdr[1] to be \"c\"")
	}

}

func TestLeaf(t *testing.T) {
	var s Sexp

	// List
	s = List{}
	if s.IsLeaf() {
		t.Error("List should never be considered a leaf, EVEN WHEN EMPTY")
	}

	s = List{List{Symbol("a"), Symbol("b")}, Symbol("c")}
	if s.LeafCount() != 3 {
		t.Error("Expected 3 leaves")
	}

	// Atom
	s = Symbol("a")
	if !s.IsLeaf() {
		t.Error("Symbols should always be leaves")
	}

	if s.LeafCount() != 1 {
		t.Error("Expected 1 leaf")
	}

	// Strict
	s = NewStrict(Symbol("a"))
	if !s.IsLeaf() {
		t.Error("Expected a leaf ")
	}

	s = addChild(s, Symbol("b"), true)
	s = addChild(s, Symbol("c"), true)

	if s.IsLeaf() {
		t.Error("This is clearly no longer a leaf node")
	}

	if s.LeafCount() != 3 {
		t.Error("Expected 3 leaves")
	}
}

// package sexp provides the data structure and parser for s-expressions in Go.
package sexp

import "fmt"

// Sexp is the interface representing a S-expression node
type Sexp interface {
	IsLeaf() bool   // indicates if this node is a leaf node
	LeafCount() int // count of leaves from the children of this node

	Head() Sexp // also known as car
	Tail() Sexp // also known as cdr

	fmt.Formatter
}

// Atom is the interface representing an atomic S-expression
type Atom interface {
	Sexp
	IsAtom() bool
}

// List is a shortcut (and default) version of a s-expression. It's used together with Symbol as its atomic s-expression
type List []Sexp

// List will always return false, even if it is empty
func (s List) IsLeaf() bool { return false }

// LeafCount returns the total number of leaves the s-expression has.
func (s List) LeafCount() (retVal int) {
	for _, child := range s {
		retVal += child.LeafCount()
	}
	return
}

func (s List) Head() Sexp { return s[0] }
func (s List) Tail() Sexp { return s[1:] }

func (s List) Format(f fmt.State, c rune) {
	fmt.Fprint(f, "(")
	for i, child := range s {
		if i == 0 {
			fmt.Fprintf(f, "%s", child)
		} else {
			fmt.Fprintf(f, " %s", child)
		}
	}
	fmt.Fprint(f, ")")
}

/* Atoms */

// Symbol is the atomic element in an s-expression.
type Symbol string

func (s Symbol) IsLeaf() bool               { return true }
func (s Symbol) LeafCount() int             { return 1 }
func (s Symbol) Head() Sexp                 { return s }
func (s Symbol) Tail() Sexp                 { return nil }
func (s Symbol) Format(f fmt.State, c rune) { fmt.Fprintf(f, "%s", string(s)) }
func (s Symbol) IsAtom() bool               { return true }

// Strict is a strict/canonical form of a s-expression, using a doubly-linked list as its backing data structure. To parse to strict s-expression, parse with strict=true
type Strict struct {
	Sexp
	parent, child Sexp
}

// NewStrict "upgrades" a Sexp into a *Strict
func NewStrict(s Sexp) *Strict {
	if ss, ok := s.(*Strict); ok {
		return ss
	}
	return &Strict{Sexp: s}
}

func (s *Strict) IsLeaf() bool { return s.child == nil }
func (s *Strict) LeafCount() int {
	if s.IsLeaf() {
		return s.Sexp.LeafCount()
	}

	count := s.Sexp.LeafCount()
	count += s.child.LeafCount()
	return count
}

func (s *Strict) Head() Sexp { return s.Sexp }
func (s Strict) Tail() Sexp  { return s.child }

func (s *Strict) Format(f fmt.State, c rune) {
	fmt.Fprintf(f, "(%s", s.Sexp)
	if s.child != nil {
		fmt.Fprintf(f, " %s", s.child)
	}
	f.Write([]byte(")"))
}

// Last is a convenience function that retuns the last value of the linked list
func (s *Strict) Last() *Strict {
	if s.child != nil {
		if child, ok := s.child.(*Strict); ok {
			return child.Last()

		}
	}
	return s
}

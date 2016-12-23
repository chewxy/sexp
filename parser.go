package sexp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type dummy struct{}

func (s dummy) IsLeaf() bool               { panic("not implemented") }
func (s dummy) LeafCount() int             { panic("not implemented") }
func (s dummy) Head() Sexp                 { panic("not implemented") }
func (s dummy) Tail() Sexp                 { panic("not implemented") }
func (s dummy) Format(f fmt.State, c rune) { fmt.Fprint(f, "DUMMY") }

const eof rune = -1

type stateFn func(*Parser) stateFn

// AtomReader is any function that takes a string and returns an Atom and an error
type AtomReader func(string) (Atom, error)

// Parser holds the state for lexing and parsing of a string into a Sexp.
type Parser struct {
	io.RuneScanner
	strict bool
	ar     AtomReader

	Output chan Sexp

	// status
	r     rune
	width int
	pos   int
	start int
	line  int
	col   int

	// state
	buf     *bytes.Buffer // the string we're reading
	stack   []Sexp        // building a stack: everytime we encounter a '(', we add a new LF. Everytime we encounter a ')', we pop the stack
	current Sexp          // current logical form
	parens  []int         // paren stack
	err     error
}

// NewParser creates a new parser from a io.Reader. Strict indicates if the string should be parsed as a S-expression with a linked-list as a backing. Optionals are AtomReader functions.
// If no AtomReader functions are passed in, the default SymbolReader will be used
func NewParser(r io.Reader, strict bool, ars ...AtomReader) *Parser {
	stack := make([]Sexp, 32) // you probably ain't gonna need more than that
	stack = stack[:0]

	var rs io.RuneScanner
	var ok bool
	if rs, ok = r.(io.RuneScanner); !ok {
		rs = bufio.NewReader(r)
	}

	var ar AtomReader
	if len(ars) != 0 {
		ar = ars[0]
	} else {
		ar = SymbolReader
	}

	return &Parser{
		RuneScanner: rs,
		strict:      strict,
		ar:          ar,

		Output: make(chan Sexp),

		width: 1,
		start: 1,
		col:   1,
		pos:   1,

		buf:   new(bytes.Buffer),
		stack: stack,
	}
}

// Run runs the parser
func (p *Parser) Run() {
	defer close(p.Output)
	for state := lexStart; state != nil; {
		state = state(p)
		if len(p.stack) == 0 && len(p.parens) == 0 && p.current != nil {
			p.Output <- p.current
			p.current = nil
		}
	}

	// for len(p.stack) > 0 {
	// 	parent := p.pop()
	// 	p.current = addChild(parent, p.current, p.strict)
	// }

	// if p.current != nil {
	// 	p.Output <- p.current
	// }
}

func (p *Parser) Error() error { return p.err }

func (p *Parser) next() rune {
	var err error
	p.r, p.width, err = p.ReadRune()
	if err == io.EOF {
		p.width = 1
		return eof
	}

	p.col += p.width
	p.pos += p.width

	return p.r
}

func (p *Parser) backup() {
	p.UnreadRune()
	p.pos -= p.width
	p.col -= p.width
}

func (p *Parser) peek() rune {
	backup := p.r
	pos := p.pos
	col := p.col

	r := p.next()
	p.backup()

	p.pos = pos
	p.col = col
	p.r = backup
	return r
}

func (p *Parser) lineCount() {
	newLines := bytes.Count(p.buf.Bytes(), []byte("\n"))

	p.line += newLines
	if newLines > 0 {
		p.col = 1
	}
}

func (p *Parser) accept() {
	p.buf.WriteRune(p.r)
}

func (p *Parser) acceptRunFn(fn func(rune) bool) {
	for fn(p.peek()) {
		p.next()
		p.accept()
	}
}

func (p *Parser) acceptRunUntilFn(fn func(rune) bool) bool {
	r := p.peek()
	for {
		if r == eof {
			return true
		}
		if !fn(r) {
			p.next()
			p.accept()
			r = p.peek()
			continue
		}
		break
	}

	return false
}

func (p *Parser) ignore() {
	p.start = p.pos
	p.buf.Reset()
}

func (p *Parser) push(s Sexp) { p.stack = append(p.stack, s) }

// func (p *Parser) pop() Sexp {
// 	if len(p.stack) == 0 {
// 		return nil
// 	}
// 	retVal := p.stack[len(p.stack)-1]
// 	p.stack = p.stack[:len(p.stack)-1]
// 	return retVal
// }

func (p *Parser) pushParens(i int) { p.parens = append(p.parens, i) }
func (p *Parser) popParens() int {
	if len(p.parens) == 0 {
		return -1
	}
	retVal := p.parens[len(p.parens)-1]
	p.parens = p.parens[:len(p.parens)-1]
	return retVal
}

func lexStart(p *Parser) stateFn {
	r := p.peek()
	switch {
	case r == '(':
		return lexParen
	case r == '#':
		return lexComment
	case r == ')':
		return lexCloseParen
	default:
		return lexSymbol
	}
}

func lexComment(p *Parser) stateFn {
	return nil
}

func lexSpace(p *Parser) stateFn {
	logf("LexSpace")
	enterLoggingContext()
	defer leaveLoggingContext()

	logf("PEEK %q", p.peek())
	p.acceptRunFn(unicode.IsSpace)
	p.lineCount()
	p.ignore()

	r := p.peek()
	switch {
	case r == '(':
		return lexParen
	case r == ')':
		return lexCloseParen
	default:
		return lexSymbol
	}

	return nil
}

func lexParen(p *Parser) stateFn {
	logf("lexParen")
	enterLoggingContext()
	defer leaveLoggingContext()

	logf("Stack %v, Current: %v | %v", p.stack, p.current, p.parens)

	p.next() // accept the paren
	p.ignore()

	if p.current != nil {
		p.push(p.current)
	}
	p.pushParens(len(p.stack))
	p.current = dummy{}

	logf("Stack %v, Current: %v | %v", p.stack, p.current, p.parens)
	return lexSpace
}

func lexSymbol(p *Parser) stateFn {
	logf("lexSymbol")
	enterLoggingContext()
	defer leaveLoggingContext()

	logf("Stack %v, Current: %v | %v", p.stack, p.current, p.parens)

	if p.acceptRunUntilFn(isSpaceOrParen) {
		return nil
	}
	word := p.buf.String()

	if p.current != nil && (p.current != dummy{}) {
		p.push(p.current)
	}

	if p.current, p.err = p.ar(word); p.err != nil {
		return nil
	}
	return lexSpace
}

func lexCloseParen(p *Parser) stateFn {
	logf("CloseParen")
	enterLoggingContext()
	defer leaveLoggingContext()

	p.next()
	p.ignore()

	parentParen := p.popParens()

	logf("Combine %v %v", p.stack[parentParen:], p.current)

	p.current = combine(p.stack[parentParen:], p.current, p.strict)
	logf("Parent Paren %v Stack %v | %v", parentParen, p.stack, p.current)
	if len(p.stack) > parentParen {
		p.stack = p.stack[:parentParen]
	}

	logf("After truncate Stack %v | %v", p.stack, p.current)

	return lexSpace
}

func addChild(parent, child Sexp, strict bool) (retVal Sexp) {
	switch p := parent.(type) {
	case Atom:
		if strict {
			childS := NewStrict(child)
			parentS := NewStrict(parent)
			childS.parent = parent
			parentS.child = child
			logf("Parent %#v", parentS)
			return parentS
		}
		retVal = List{parent, child}
		return
	case List:
		retVal = append(p, child)
		return
	case *Strict:
		if s, ok := child.(*Strict); ok {
			s.parent = parent
		} else {
			child = &Strict{Sexp: child, parent: parent}
		}

		p.Last().child = child
		return p
	}
	panic("Unreachable")
}

func combine(sexps []Sexp, cur Sexp, strict bool) Sexp {
	logf("Combinnig")
	enterLoggingContext()
	defer leaveLoggingContext()

	switch len(sexps) {
	case 0:
		return cur
	case 1:
		sexps[0] = addChild(sexps[0], cur, strict)
		return sexps[0]
	default:
		if strict {
			for i := len(sexps) - 1; i >= 0; i-- {
				cur = addChild(sexps[i], cur, strict)
			}
			return cur
		}

		for _, s := range sexps[1:] {
			sexps[0] = addChild(sexps[0], s, strict)
		}
		sexps[0] = addChild(sexps[0], cur, strict)
		return sexps[0]
	}
}

func isSpaceOrParen(r rune) bool { return r == '(' || r == ')' || unicode.IsSpace(r) }

// SymbolReader is the default AtomReader.
func SymbolReader(s string) (Atom, error) { return Symbol(s), nil }

package sexp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var examples = []string{
	//  basics
	"(lambda x (f x))",

	`(defun factorial (x)
	 		(if (zerop x)
	     		1
	     		(* x (factorial (- x 1)))))
	`,

	"(lambda x (f x)) (lambda x (f x))",
}

var corrects = []string{
	"(lambda x (f x))",
	"(defun factorial x (if (zerop x) 1 (* x (factorial (- x 1)))))",
	"(lambda x (f x)) (lambda x (f x))",
}

func TestParse(t *testing.T) {
	var buf bytes.Buffer
	for i, ex := range examples {
		p := NewParser(strings.NewReader(ex), false)

		done := make(chan struct{})
		var sexps []Sexp
		go func(ch chan Sexp, done chan struct{}) {
			for s := range ch {
				sexps = append(sexps, s)
			}

			done <- struct{}{}
		}(p.Output, done)

		p.Run()
		<-done

		if p.Error() != nil {
			t.Error(p.err)
		}

		if len(sexps) == 0 {
			t.Fatal("Expected more than 0 sexps")
		}

		for i, s := range sexps {
			if i < len(sexps)-1 {
				fmt.Fprintf(&buf, "%v ", s)
			} else {
				fmt.Fprintf(&buf, "%v", s)
			}
		}

		s := buf.String()

		if s != corrects[i] {
			t.Errorf("Example %d Expected %q. Got %q", i, corrects[i], s)
		}
		buf.Reset()

		// if fmt.Sprintf("%s", p.current) != corrects[i] {
		//
		// }
	}
}

var strictCorrects = []string{
	"(lambda (x (f x)))",

	"(defun (factorial (x (if (zerop (1 (* (x (factorial (- (x 1)))))))))))",
	"(lambda (x (f x))) (lambda (x (f x)))",
}

func TestParseStrict(t *testing.T) {
	var buf bytes.Buffer
	for i, ex := range examples {
		p := NewParser(strings.NewReader(ex), true)

		done := make(chan struct{})
		var sexps []Sexp
		go func(ch chan Sexp, done chan struct{}) {
			for s := range ch {
				sexps = append(sexps, s)
			}

			done <- struct{}{}
		}(p.Output, done)

		p.Run()
		<-done

		if p.Error() != nil {
			t.Error(p.err)
		}

		if len(sexps) == 0 {
			t.Fatal("Expected more than 0 sexps")
		}

		for i, s := range sexps {
			if i < len(sexps)-1 {
				fmt.Fprintf(&buf, "%v ", s)
			} else {
				fmt.Fprintf(&buf, "%v", s)
			}
		}

		s := buf.String()

		if s != strictCorrects[i] {
			t.Errorf("Example %d Expected %q. Got %q", i, strictCorrects[i], s)
		}
		buf.Reset()
	}
}

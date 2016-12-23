# sexp [![Build Status](https://travis-ci.org/chewxy/sexp.svg?branch=master)](https://travis-ci.org/chewxy/sexp) [![Coverage Status](https://coveralls.io/repos/github/chewxy/sexp/badge.svg?branch=master)](https://coveralls.io/github/chewxy/sexp?branch=master) [![GoDoc](https://godoc.org/github.com/chewxy/sexp?status.svg)](https://godoc.org/github.com/chewxy/sexp)
package sexp provides the data structure and parser for [s-expressions](https://en.wikipedia.org/wiki/S-expression) in Go.

# Installation

`go get -u github.com/chewxy/sexp`

This package does not use any dependencies other than those found in the standard library.

# FAQ

### 1. Why another package?
Why write another package for s-expressions when there are two other great packages - one [by Shane Hanna](https://github.com/shanna/sexp) and another [by nsf](https://github.com/nsf/sexp)? In fact, s-expressions are so simple, they can be snippetted in [Rosetta Code](https://rosettacode.org/wiki/S-Expressions). So why yet another package?

The main reason is simple - I work with a lot of s-expressions in different formats. I wanted the ability to reuse a lot of parsing code to parse into different types of `Atom`s and `Sexp`s. More importantly, I needed the ability to compose my parsers. Making a general parser that implements `io.Runescanner` means this `*Parser` is composable. 

### 2. Why `AtomReader` instead of making `Sexp`s implement `MarshalText` and `UnmarshalText`?##
The reason for this is mainly historical. I have a large number of libraries which use the `AtomReader` format instead of `MarshalText`. I'd argue that `MarshalText` and `UnmarshalText` is more elegant, but we're all stuck with the cruft of our past, eh?

### 3. Your definition of "Reader" is weird. 
Yes. It is. `AtomReader` reads a string and returns a `Atom`. It doesn't read the same way as `io.Reader` does. It's retained the `Reader` suffix because of historical reasons. As an aside, wouldn't it be nice if we could do this: `go:generate reverseStringer --type=...`, kinda like a `deriving Read` in  Haskell.

### 4. Linecounts and  stuff would be nice, especially for error reporting
Yes. It would be. I started writing it with those in mind but ended up not having them due to time constraints. Never really went back to adding it. Feel free to send a pull request for that.

### 5. Comments are not supported??!!!
Yes. For now. The basics for handling comments are in the code. I just never finished it. Feel free to send a pull request.

# Contributing 

1. File an issue on github if you find any.
2. Fork the project on Github.
3. Write your changes.
4. Make sure the tests (which are sparse for now) passes.
5. Send a pull request. 
6. I'll review the code, suggest changes if any, and then merge.

Topic branches aren't necessary because the package is small. Most changes will be directly merged into Master.

# Licence
Package sexp is licenced under the MIT licence.
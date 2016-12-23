// +build !debug

package sexp

// DEBUG indicates whether this package is in debug mode.
const DEBUG = false

var _TABCOUNT uint32

func tabcount() int                             { return -1 }
func enterLoggingContext()                      {}
func leaveLoggingContext()                      {}
func logf(format string, others ...interface{}) {}

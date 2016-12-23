// +build !debug

package sexp

const DEBUG = false

var _TABCOUNT uint32

func tabcount() int                             { return -1 }
func enterLoggingContext()                      {}
func leaveLoggingContext()                      {}
func logf(format string, others ...interface{}) {}

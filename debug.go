// +build debug

package sexp

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
)

const DEBUG = false

var _TABCOUNT uint32

var TRACK = false

var _logger_ = log.New(os.Stderr, "", 0)
var replacement = "\n"

func tabcount() int {
	return int(atomic.LoadUint32(&_TABCOUNT))
}

func enterLoggingContext() {
	atomic.AddUint32(&_TABCOUNT, 1)
	tabcount := tabcount()
	_logger_.SetPrefix(strings.Repeat("\t", tabcount))
	replacement = "\n" + strings.Repeat("\t", tabcount)
}

func leaveLoggingContext() {
	tabcount := tabcount()
	tabcount--

	if tabcount < 0 {
		atomic.StoreUint32(&_TABCOUNT, 0)
		tabcount = 0
	} else {
		atomic.StoreUint32(&_TABCOUNT, uint32(tabcount))
	}
	_logger_.SetPrefix(strings.Repeat("\t", tabcount))
	replacement = "\n" + strings.Repeat("\t", tabcount)
}

func logf(format string, others ...interface{}) {
	if DEBUG {
		// format = strings.Replace(format, "\n", replacement, -1)
		s := fmt.Sprintf(format, others...)
		s = strings.Replace(s, "\n", replacement, -1)
		_logger_.Println(s)
		// _logger_.Printf(format, others...)
	}
}

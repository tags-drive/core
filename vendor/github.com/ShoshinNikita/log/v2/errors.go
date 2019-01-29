package clog

import (
	"fmt"
	"time"
)

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Error(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getErrMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Errorf(format string, v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getErrMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())
}

// Errorln prints error
// Output pattern: (?time) [ERR] (?file:line) error
func (l Logger) Errorln(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getErrMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

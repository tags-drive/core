package clog

import (
	"fmt"
	"time"
)

// Warn prints warning
// Output pattern: (?time) [WRN] warning
func (l Logger) Warn(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getWarnMsg())
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

// Warnf prints warning
// Output pattern: (?time) [WRN] warning
func (l Logger) Warnf(format string, v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getWarnMsg())
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())
}

// Warnln prints warning
// Output pattern: (?time) [WRN] warning
func (l Logger) Warnln(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getWarnMsg())
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

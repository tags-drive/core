package clog

import (
	"fmt"
	"time"
)

// Print prints msg
// Output pattern: (?time) msg
func (l Logger) Print(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

// Printf prints msg
// Output pattern: (?time) msg
func (l Logger) Printf(format string, v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())
}

// Println prints msg
// Output pattern: (?time) msg
func (l Logger) Println(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

package clog

import (
	"fmt"
	"time"
)

// Info prints info message
// Output pattern: (?time) [INF] msg
func (l Logger) Info(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getInfoMsg())
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

// Infof prints info message
// Output pattern: (?time) [INF] msg
func (l Logger) Infof(format string, v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getInfoMsg())
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())
}

// Infoln prints info message
// Output pattern: (?time) [INF] msg
func (l Logger) Infoln(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getInfoMsg())
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

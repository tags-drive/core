package clog

import (
	"fmt"
	"time"
)

// Debug prints debug message if Debug mode is on
// Output pattern: (?time) [DBG] msg
func (l Logger) Debug(v ...interface{}) {
	if !l.debug {
		return
	}

	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getDebugMsg())
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

// Debugf prints debug message if Debug mode is on
// Output pattern: (?time) [DBG] msg
func (l Logger) Debugf(format string, v ...interface{}) {
	if !l.debug {
		return
	}

	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getDebugMsg())
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())
}

// Debugln prints debug message if Debug mode is on
// Output pattern: (?time) [DBG] msg
func (l Logger) Debugln(v ...interface{}) {
	if !l.debug {
		return
	}

	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getDebugMsg())
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())
}

package clog

import (
	"fmt"
	"os"
	"time"
)

// Fatal prints error and call os.Exit(1)
// Output pattern: (?time) [FAT] (?file:line) error
func (l Logger) Fatal(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getFatalMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprint(l.buff, v...)

	l.output.Write(l.buff.Bytes())

	os.Exit(1)
}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FAT] (?file:line) error
func (l Logger) Fatalf(format string, v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getFatalMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprintf(l.buff, format, v...)

	l.output.Write(l.buff.Bytes())

	os.Exit(1)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FAT] (?file:line) error
func (l Logger) Fatalln(v ...interface{}) {
	now := time.Now()

	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.buff.Reset()

	l.buff.Write(l.getTime(now))
	l.buff.Write(l.getFatalMsg())
	l.buff.Write(l.getCaller())
	fmt.Fprintln(l.buff, v...)

	l.output.Write(l.buff.Bytes())

	os.Exit(1)
}

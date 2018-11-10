package log

import (
	"fmt"
	"os"
)

// Fatal prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatal(v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getFatalMsg()
	if l.printErrorLine {
		text += l.getCaller()
	}
	l.printText(text + fmt.Sprint(v...))
	os.Exit(1)

}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatalf(format string, v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getFatalMsg()
	if l.printErrorLine {
		text += l.getCaller()
	}
	l.printText(text + fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func (l Logger) Fatalln(v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getFatalMsg()
	if l.printErrorLine {
		text += l.getCaller()
	}
	l.printText(text + fmt.Sprint(v...))
	os.Exit(1)
}

package log

import (
	"fmt"
)

// Info prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Info(v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getInfoMsg()
	l.printText(text + fmt.Sprint(v...))
}

// Infof prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Infof(format string, v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getInfoMsg()
	l.printText(text + fmt.Sprintf(format, v...))
}

// Infoln prints info message
// Output pattern: (?time) [INFO] msg
func (l Logger) Infoln(v ...interface{}) {
	text := ""
	if l.printTime {
		text = l.getTime()
	}
	text += l.getInfoMsg()
	l.printText(text + fmt.Sprintln(v...))
}

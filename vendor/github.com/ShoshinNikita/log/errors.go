package log

import (
	"fmt"
)

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Error(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getErrMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprint(v...))
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorf(format string, v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getErrMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprintf(format, v...))
}

// Errorln prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorln(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getErrMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprintln(v...))
}

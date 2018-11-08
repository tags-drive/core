package log

import (
	"fmt"
)

// Warn prints warning
// Output pattern: (?time) [WARN] warning
func Warn(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getWarnMsg()
	printText(text + fmt.Sprint(v...))
}

// Warnf prints warning
// Output pattern: (?time) [WARN] warning
func Warnf(format string, v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getWarnMsg()
	printText(text + fmt.Sprintf(format, v...))
}

// Warnln prints warning
// Output pattern: (?time) [WARN] warning
func Warnln(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getWarnMsg()
	printText(text + fmt.Sprintln(v...))
}

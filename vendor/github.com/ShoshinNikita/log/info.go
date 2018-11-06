package log

import (
	"fmt"
)

func Info(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getInfoMsg()
	printText(text + fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getInfoMsg()
	printText(text + fmt.Sprintf(format, v...))
}

func Infoln(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getInfoMsg()
	printText(text + fmt.Sprintln(v...))
}

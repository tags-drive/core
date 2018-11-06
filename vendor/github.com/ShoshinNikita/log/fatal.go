package log

import (
	"fmt"
	"os"
)

func Fatal(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getFatalMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprint(v...))
	os.Exit(1)

}

func Fatalf(format string, v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getFatalMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatalln(v ...interface{}) {
	text := ""
	if printTime {
		text = getTime()
	}
	text += getFatalMsg()
	if printErrorLine {
		text += getCaller()
	}
	printText(text + fmt.Sprint(v...))
	os.Exit(1)
}

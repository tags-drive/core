package log

import (
	"fmt"
)

func Print(v ...interface{}) {
	printText(fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	printText(fmt.Sprintf(format, v...))
}

func Println(v ...interface{}) {
	printText(fmt.Sprintln(v...))
}

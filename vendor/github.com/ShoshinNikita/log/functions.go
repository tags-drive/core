package log

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
)

var (
	// For time
	timePrintf = color.New(color.FgHiGreen).SprintfFunc()

	// For [INFO]
	infoPrint = color.New(color.FgCyan).SprintFunc()

	// For [WARN]
	warnPrint = color.New(color.FgYellow).SprintFunc()

	// For [ERR]
	errorPrint = color.New(color.FgRed).SprintFunc()

	// For [FATAL]
	fatalPrint = color.New(color.BgRed).SprintFunc()
)

func (l Logger) getCaller() string {
	var (
		file string
		line int
		ok   bool
	)

	if l.global {
		_, file, line, ok = runtime.Caller(3)
	} else {
		_, file, line, ok = runtime.Caller(2)
	}
	if !ok {
		return ""
	}

	var shortFile string
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			shortFile = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d ", shortFile, line)
}

func (l Logger) getTime() string {
	if l.printColor {
		return timePrintf("%s ", time.Now().Format(timeLayout))
	}
	return fmt.Sprintf("%s ", time.Now().Format(timeLayout))
}

func (l Logger) getInfoMsg() string {
	if l.printColor {
		return infoPrint("[INFO] ")
	}
	return "[INFO] "
}

func (l Logger) getWarnMsg() string {
	if l.printColor {
		return warnPrint("[WARN] ")
	}
	return "[WARN] "
}

func (l Logger) getErrMsg() string {
	if l.printColor {
		return errorPrint("[ERR] ")
	}
	return "[ERR] "
}

func (l Logger) getFatalMsg() (s string) {
	if l.printColor {
		return fatalPrint("[FATAL]") + " "
	}
	return "[FATAL] "
}

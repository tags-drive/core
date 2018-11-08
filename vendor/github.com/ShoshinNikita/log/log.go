// Package log provides functions for pretty print
//
// Patterns of functions print:
// * Print(), Printf(), Println():
//   (?time) msg
// * Info(), Infof(), Infoln():
//   (?time) [INFO] msg
// * Error(), Errorf(), Errorln():
//   (?time) [ERR] (?file:line) error
// * Fatal(), Fatalf(), Fatalln():
//   (?time) [FATAL] (?file:line) error
//
// Time pattern: MM.dd.yyyy hh:mm:ss (01.30.2018 05:5:59)
//
package log

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
)

type textStruct struct {
	text string
	ch   chan struct{}
}

func (t *textStruct) done() {
	close(t.ch)
}

func newText(text string) textStruct {
	return textStruct{text: text, ch: make(chan struct{})}
}

const (
	timeLayout = "01.02.2006 15:04:05"
)

var (
	printTime      bool
	printColor     = true
	printErrorLine = true

	printChan = make(chan textStruct, 500)

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

// init runs goroutine, which prints text from channel
func init() {
	go func() {
		for text := range printChan {
			fmt.Fprint(color.Output, text.text)
			text.done()
		}
	}()
}

func printText(text string) {
	t := newText(text)
	printChan <- t
	<-t.ch
}

// PrintTime sets printTime
// Time isn't printed by default
func PrintTime(b bool) {
	printTime = b
}

// ShowTime sets printTime
// Time isn't printed by default
//
// It was left for backwards compatibility
var ShowTime = PrintTime

// PrintColor sets printColor
// printColor is true by default
func PrintColor(b bool) {
	printColor = b
}

// PrintErrorLine sets PrintErrorLine
// If PrintErrorLine is true, log.Error(), log.Errorf(), log.Errorln() will print file and line,
// where functions were called.
// PrintErrorLine is true by default
func PrintErrorLine(b bool) {
	printErrorLine = b
}

func getCaller() string {
	// We need to skip 2 functions (this and log.Error(), log.Errorf() or log.Errorln())
	_, file, line, ok := runtime.Caller(2)
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

func getTime() string {
	if printColor {
		return timePrintf("%s ", time.Now().Format(timeLayout))
	}
	return fmt.Sprintf("%s ", time.Now().Format(timeLayout))
}

func getInfoMsg() string {
	if printColor {
		return infoPrint("[INFO] ")
	}
	return "[INFO] "
}

func getWarnMsg() string {
	if printColor {
		return warnPrint("[WARN] ")
	}
	return "[WARN] "
}

func getErrMsg() string {
	if printColor {
		return errorPrint("[ERR] ")
	}
	return "[ERR] "
}

func getFatalMsg() (s string) {
	if printColor {
		return fatalPrint("[FATAL]") + " "
	}
	return "[FATAL] "
}

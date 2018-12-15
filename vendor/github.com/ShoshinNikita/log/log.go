// Package log provides functions for pretty print
//
// Patterns of functions print:
// * Print(), Printf(), Println():
//   (?time) msg
// * Info(), Infof(), Infoln():
//   (?time) [INFO] msg
// * Warn(), Warnf(), Warnln():
//   (?time) [WARN] warning
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

	"github.com/fatih/color"
)

const (
	timeLayout = "01.02.2006 15:04:05"
)

// init inits globalLogger with NewLogger()
func init() {
	globalLogger = NewLogger()

	globalLogger.PrintTime(false)
	globalLogger.PrintColor(false)
	globalLogger.PrintErrorLine(false)

	globalLogger.global = true
}

type textStruct struct {
	text string
	ch   chan struct{}
}

func newText(text string) textStruct {
	return textStruct{text: text, ch: make(chan struct{})}
}

func (t *textStruct) done() {
	close(t.ch)
}

type Logger struct {
	printTime      bool
	printColor     bool
	printErrorLine bool

	printChan chan textStruct
	global    bool
}

// NewLogger creates *Logger and run goroutine (Logger.printer())
func NewLogger() *Logger {
	l := new(Logger)
	l.printChan = make(chan textStruct)
	go l.printer()
	return l
}

func (l Logger) printer() {
	for text := range l.printChan {
		fmt.Fprint(color.Output, text.text)
		text.done()
	}
}

func (l Logger) printText(text string) {
	t := newText(text)
	l.printChan <- t
	<-t.ch
}

// PrintTime sets Logger.printTime to b
func (l *Logger) PrintTime(b bool) {
	l.printTime = b
}

// PrintColor sets Logger.printColor to b
func (l *Logger) PrintColor(b bool) {
	l.printColor = b
}

// PrintErrorLine sets Logger.printErrorLine to b
func (l *Logger) PrintErrorLine(b bool) {
	l.printErrorLine = b
}

var globalLogger *Logger

// PrintTime sets globalLogger.PrintTime
// Time isn't printed by default
func PrintTime(b bool) {
	globalLogger.PrintTime(b)
}

// ShowTime sets printTime
// Time isn't printed by default
//
// It was left for backwards compatibility
var ShowTime = PrintTime

// PrintColor sets printColor
// printColor is false by default
func PrintColor(b bool) {
	globalLogger.PrintColor(b)
}

// PrintErrorLine sets PrintErrorLine
// If PrintErrorLine is true, log.Error(), log.Errorf(), log.Errorln() will print file and line,
// where functions were called.
// PrintErrorLine is false by default
func PrintErrorLine(b bool) {
	globalLogger.PrintErrorLine(b)
}

/* Print */

// Print prints msg
// Output pattern: (?time) msg
func Print(v ...interface{}) {
	globalLogger.Print(v...)
}

// Printf prints msg
// Output pattern: (?time) msg
func Printf(format string, v ...interface{}) {
	globalLogger.Printf(format, v...)
}

// Println prints msg
// Output pattern: (?time) msg
func Println(v ...interface{}) {
	globalLogger.Println(v...)
}

/* Info */

// Info prints info message
// Output pattern: (?time) [INFO] msg
func Info(v ...interface{}) {
	globalLogger.Info(v...)
}

// Infof prints info message
// Output pattern: (?time) [INFO] msg
func Infof(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
}

// Infoln prints info message
// Output pattern: (?time) [INFO] msg
func Infoln(v ...interface{}) {
	globalLogger.Infoln(v...)
}

/* Warn */

// Warn prints warning
// Output pattern: (?time) [WARN] warning
func Warn(v ...interface{}) {
	globalLogger.Warn(v...)
}

// Warnf prints warning
// Output pattern: (?time) [WARN] warning
func Warnf(format string, v ...interface{}) {
	globalLogger.Warnf(format, v...)
}

// Warnln prints warning
// Output pattern: (?time) [WARN] warning
func Warnln(v ...interface{}) {
	globalLogger.Warnln(v...)
}

/* Error */

// Error prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Error(v ...interface{}) {
	globalLogger.Error(v...)
}

// Errorf prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorf(format string, v ...interface{}) {
	globalLogger.Errorf(format, v...)
}

// Errorln prints error
// Output pattern: (?time) [ERR] (?file:line) error
func Errorln(v ...interface{}) {
	globalLogger.Errorln(v...)
}

/* Fatal */

// Fatal prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatal(v ...interface{}) {
	globalLogger.Fatal(v...)
}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FATAL] (?file:line) error
func Fatalln(v ...interface{}) {
	globalLogger.Fatalln(v...)
}

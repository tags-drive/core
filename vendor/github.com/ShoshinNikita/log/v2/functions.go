package clog

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
)

var (
	usualDEBUG = []byte("[DBG] ")
	usualINFO  = []byte("[INF] ")
	usualWARN  = []byte("[WRN] ")
	usualERR   = []byte("[ERR] ")
	usualFATAL = []byte("[FAT] ")

	coloredDEBUG = []byte(color.HiMagentaString(string(usualDEBUG)))
	coloredINFO  = []byte(color.CyanString(string(usualINFO)))
	coloredWARN  = []byte(color.YellowString(string(usualWARN)))
	coloredERR   = []byte(color.RedString(string(usualERR)))
	coloredFATAL = []byte(color.New(color.BgRed).Sprint("[FAT]") + " ")

	timePrintf   = color.New(color.FgHiGreen).SprintfFunc()
	callerPrintf = color.RedString // color is the same as coloredErr
)

// getTime returns "file:line" if l.printErrorLine == true, else it returns empty string
func (l Logger) getCaller() []byte {
	if !l.printErrorLine {
		return nil
	}

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
		return nil
	}

	var shortFile string
	for i := len(file) - 1; i >= 0; i-- {
		if file[i] == '/' {
			shortFile = file[i+1:]
			break
		}
	}

	if l.printColor {
		return []byte(callerPrintf("%s:%d ", shortFile, line))
	}
	return []byte(fmt.Sprintf("%s:%d ", shortFile, line))
}

// getTime returns time if l.printTime == true, else it returns empty string
func (l Logger) getTime(t time.Time) []byte {
	if !l.printTime {
		return nil
	}

	if l.printColor {
		return []byte(timePrintf("%s ", t.Format(l.timeLayout)))
	}
	return []byte(t.Format(l.timeLayout) + " ")
}

func (l Logger) getDebugMsg() []byte {
	if l.printColor {
		return coloredDEBUG
	}
	return usualDEBUG
}

func (l Logger) getInfoMsg() []byte {
	if l.printColor {
		return coloredINFO
	}
	return usualINFO
}

func (l Logger) getWarnMsg() []byte {
	if l.printColor {
		return coloredWARN
	}
	return usualWARN
}

func (l Logger) getErrMsg() []byte {
	if l.printColor {
		return coloredERR
	}
	return usualERR
}

func (l Logger) getFatalMsg() []byte {
	if l.printColor {
		return coloredFATAL
	}
	return usualFATAL
}

// Package log provides functions for pretty print
//
// Patterns of functions print:
//
// * (?time) msg - Print(), Printf(), Println():
//
// * (?time) [DBG] msg - Debug(), Debugf(), Debugln():
//
// * (?time) [INF] msg - Info(), Infof(), Infoln():
//
// * (?time) [WRN] warning - Warn(), Warnf(), Warnln():
//
// * (?time) [ERR] (?file:line) error - Error(), Errorf(), Errorln():
//
// * (?time) [FAT] (?file:line) error - Fatal(), Fatalf(), Fatalln():
//
// Time pattern: MM.dd.yyyy hh:mm:ss (01.30.2018 05:5:59)
//
package clog

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/fatih/color"
)

const (
	DefaultTimeLayout = "01.02.2006 15:04:05"
)

type Logger struct {
	output io.Writer
	mutex  *sync.Mutex
	buff   *bytes.Buffer

	global bool

	debug          bool
	printTime      bool
	printColor     bool
	printErrorLine bool
	timeLayout     string
}

func (l *Logger) Write(b []byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.output.Write(b)
}

func (l *Logger) WriteString(s string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.output.Write([]byte(s))
}

func NewDevLogger() *Logger {
	return NewDevConfig().Build()
}

func NewProdLogger() *Logger {
	return NewProdConfig().Build()
}

type Config struct {
	output         io.Writer
	debug          bool
	printTime      bool
	printColor     bool
	printErrorLine bool
	timeLayout     string
}

func NewDevConfig() *Config {
	return &Config{
		output:         color.Output,
		debug:          true,
		printTime:      true,
		printColor:     true,
		printErrorLine: true,
		timeLayout:     DefaultTimeLayout,
	}
}

func NewProdConfig() *Config {
	return &Config{
		output:         os.Stdout,
		debug:          false,
		printTime:      true,
		printColor:     false,
		printErrorLine: true,
		timeLayout:     DefaultTimeLayout,
	}
}

// Build create a new Logger according to Config
func (c *Config) Build() *Logger {
	l := new(Logger)
	l.mutex = new(sync.Mutex)
	l.buff = new(bytes.Buffer)

	switch {
	case c.printColor && c.output == nil:
		l.output = color.Output
	case c.output != nil:
		l.output = c.output
	default:
		l.output = os.Stdout
	}

	l.debug = c.debug
	l.printTime = c.printTime
	l.printColor = c.printColor
	l.printErrorLine = c.printErrorLine

	l.timeLayout = DefaultTimeLayout
	if c.timeLayout != "" {
		l.timeLayout = c.timeLayout
	}

	return l
}

// Debug sets Config.debug to b
func (c *Config) Debug(b bool) *Config {
	c.debug = b
	return c
}

// PrintTime sets Config.printTime to b
func (c *Config) PrintTime(b bool) *Config {
	c.printTime = b
	return c
}

// PrintColor sets Config.printColor to b
func (c *Config) PrintColor(b bool) *Config {
	c.printColor = b
	return c
}

// PrintErrorLine sets Config.printErrorLine to b
func (c *Config) PrintErrorLine(b bool) *Config {
	c.printErrorLine = b
	return c
}

// SetOutput changes Config.output writer.
func (c *Config) SetOutput(w io.Writer) *Config {
	c.output = w
	return c
}

// SetTimeLayout changes Config.timeLayout
// Default Config.timeLayout is DefaultTimeLayout
func (c *Config) SetTimeLayout(layout string) *Config {
	c.timeLayout = layout
	return c
}

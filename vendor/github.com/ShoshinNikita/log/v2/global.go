package clog

var globalLogger *Logger

// init inits globalLogger with NewLogger()
func init() {
	globalLogger = NewProdLogger()
	globalLogger.global = true
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
// Output pattern: (?time) [INF] msg
func Info(v ...interface{}) {
	globalLogger.Info(v...)
}

// Infof prints info message
// Output pattern: (?time) [INF] msg
func Infof(format string, v ...interface{}) {
	globalLogger.Infof(format, v...)
}

// Infoln prints info message
// Output pattern: (?time) [INF] msg
func Infoln(v ...interface{}) {
	globalLogger.Infoln(v...)
}

/* Warn */

// Warn prints warning
// Output pattern: (?time) [WRN] warning
func Warn(v ...interface{}) {
	globalLogger.Warn(v...)
}

// Warnf prints warning
// Output pattern: (?time) [WRN] warning
func Warnf(format string, v ...interface{}) {
	globalLogger.Warnf(format, v...)
}

// Warnln prints warning
// Output pattern: (?time) [WRN] warning
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
// Output pattern: (?time) [FAT] (?file:line) error
func Fatal(v ...interface{}) {
	globalLogger.Fatal(v...)
}

// Fatalf prints error and call os.Exit(1)
// Output pattern: (?time) [FAT] (?file:line) error
func Fatalf(format string, v ...interface{}) {
	globalLogger.Fatalf(format, v...)
}

// Fatalln prints error and call os.Exit(1)
// Output pattern: (?time) [FAT] (?file:line) error
func Fatalln(v ...interface{}) {
	globalLogger.Fatalln(v...)
}

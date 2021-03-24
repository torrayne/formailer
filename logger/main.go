package logger

import (
	"log"
	"os"
)

const logFlags = log.LstdFlags | log.LUTC | log.Lmsgprefix

var logInfo = log.New(os.Stdout, "info: ", logFlags)
var logError = log.New(os.Stderr, "error: ", logFlags)

// Info logs to os.Stdout using println prefixed with info:
func Info(v ...interface{}) {
	logInfo.Println(v...)
}

// Infof logs to os.Stdout prefixed using printf with info:
func Infof(format string, v ...interface{}) {
	logInfo.Printf(format, v...)
}

// Error logs to os.Stderr using println prefixed with error:
func Error(v ...interface{}) {
	logError.Println(v...)
}

// Errorf logs to os.Stderr prefixed using printf with error:
func Errorf(format string, v ...interface{}) {
	logError.Printf(format, v...)
}

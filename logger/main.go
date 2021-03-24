package logger

import (
	"log"
	"os"
)

const logFlags = log.LstdFlags | log.LUTC | log.Lmsgprefix

var logInfo = log.New(os.Stdout, "info: ", logFlags)
var logError = log.New(os.Stderr, "error: ", logFlags)

func Info(v ...interface{}) {
	logInfo.Println(v...)
}

func Infof(format string, v ...interface{}) {
	logInfo.Printf(format, v...)
}

func Error(v ...interface{}) {
	logError.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	logError.Printf(format, v...)
}

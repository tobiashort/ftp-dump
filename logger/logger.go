package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	flags       = log.Ldate | log.Ltime | log.Lshortfile
	fatalLogger = log.New(os.Stdout, "[FATAL] ", flags)
	errorLogger = log.New(os.Stdout, "[ERROR] ", flags)
	infoLogger  = log.New(os.Stdout, "[INFO]  ", flags)
)

func Fatal(v ...interface{}) {
	fatalLogger.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Error(v ...interface{}) {
	errorLogger.Output(2, fmt.Sprintln(v...))
}

func Info(v ...interface{}) {
	infoLogger.Output(2, fmt.Sprintln(v...))
}

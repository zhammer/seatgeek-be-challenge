package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Logger struct {
	debug  bool
	logger *log.Logger
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.log("DEBUG", format, v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log("INFO ", format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log("ERROR", format, v...)
}

func (l *Logger) log(level string, format string, v ...interface{}) {
	mask := fmt.Sprintf("[%s] - %s - %s", logPrefix(), level, format)
	l.logger.Printf(mask, v...)
}

func NewLogger(debug bool) *Logger {
	actualLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	actualLogger.Printf("debug enabled: [%s]", strconv.FormatBool(debug))
	return &Logger{
		debug:  debug,
		logger: actualLogger,
	}
}

package logger

import (
	"log"
	"os"
)

type Level int

const (
	ERROR Level = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var defaultLevel = WARN

type Logger struct {
	level  Level
	logger *log.Logger
}

func NewLogger(level Level) *Logger {
	return &Logger{level, log.New(os.Stderr, "", log.LstdFlags)}
}

func NewLoggerWithStringLogLevel(s string) *Logger {
	var level Level
	switch s {
	case "ERROR":
		level = ERROR
	case "WARN":
		level = WARN
	case "INFO":
		level = INFO
	case "DEBUG":
		level = DEBUG
	case "TRACE":
		level = TRACE
	default:
		level = defaultLevel
	}
	return NewLogger(level)
}

func (l Logger) Trace(m string, f ...interface{}) {
	if l.level < TRACE {
		return
	}
	l.logger.Printf(m, f...)
}

func (l Logger) Debug(m string, f ...interface{}) {
	if l.level < DEBUG {
		return
	}
	l.logger.Printf(m, f...)
}

func (l Logger) Info(m string, f ...interface{}) {
	if l.level < INFO {
		return
	}
	l.logger.Printf(m, f...)
}

func (l Logger) Warn(m string, f ...interface{}) {
	if l.level < WARN {
		return
	}
	l.logger.Printf(m, f...)
}

func (l Logger) Error(m string, f ...interface{}) {
	if l.level < ERROR {
		return
	}
	l.logger.Printf(m, f...)
}

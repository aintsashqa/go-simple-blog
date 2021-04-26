package logger

import (
	"io"
)

const (
	InfoLogLevel LogLevel = iota
	ErrorLogLevel
	CriticalLogLevel
)

type LogLevel uint8

type Logger interface {
	Print(LogLevel, ...interface{})
	Printf(LogLevel, string, ...interface{})

	Info(...interface{})
	Infof(string, ...interface{})

	Error(...interface{})
	Errorf(string, ...interface{})

	Critical(...interface{})
	Criticalf(string, ...interface{})
}

type LoggerProvider interface {
	Logger

	SetOutput(io.Writer)
	SetOutputs(...io.Writer)
}

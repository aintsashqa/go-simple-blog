package standart

import (
	"io"
	"log"

	"github.com/aintsashqa/go-simple-blog/pkg/logger"
)

type StandartLoggerProvider struct{}

func NewStandartLoggerProvider() *StandartLoggerProvider {
	return &StandartLoggerProvider{}
}

func (l *StandartLoggerProvider) Print(level logger.LogLevel, a ...interface{}) {
	switch level {
	case logger.InfoLogLevel:
		l.Info(a...)

	case logger.ErrorLogLevel:
		l.Error(a...)

	case logger.CriticalLogLevel:
		l.Critical(a...)
	}
}

func (l *StandartLoggerProvider) Printf(level logger.LogLevel, f string, a ...interface{}) {
	switch level {
	case logger.InfoLogLevel:
		l.Infof(f, a...)

	case logger.ErrorLogLevel:
		l.Errorf(f, a...)

	case logger.CriticalLogLevel:
		l.Criticalf(f, a...)
	}
}

func (l *StandartLoggerProvider) Info(a ...interface{}) {
	log.SetPrefix("[INFO] ")
	log.Print(a...)
}

func (l *StandartLoggerProvider) Infof(f string, a ...interface{}) {
	log.SetPrefix("[INFO] ")
	log.Printf(f, a...)
}

func (l *StandartLoggerProvider) Error(a ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Print(a...)
}

func (l *StandartLoggerProvider) Errorf(f string, a ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Printf(f, a...)
}

func (l *StandartLoggerProvider) Critical(a ...interface{}) {
	log.SetPrefix("[CRITICAL] ")
	log.Fatal(a...)
}

func (l *StandartLoggerProvider) Criticalf(f string, a ...interface{}) {
	log.SetPrefix("[CRITICAL] ")
	log.Fatalf(f, a...)
}

func (l *StandartLoggerProvider) SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func (l *StandartLoggerProvider) SetOutputs(ws ...io.Writer) {
	w := io.MultiWriter(ws...)
	l.SetOutput(w)
}

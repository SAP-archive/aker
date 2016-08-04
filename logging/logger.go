package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Logger interface {
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
}

var DefaultLogger Logger = NewNativeLogger(os.Stdout, os.Stderr)

func Debugf(message string, args ...interface{}) {
	DefaultLogger.Debugf(message, args...)
}

func Infof(message string, args ...interface{}) {
	DefaultLogger.Infof(message, args...)
}

func Warnf(message string, args ...interface{}) {
	DefaultLogger.Warnf(message, args...)
}

func Errorf(message string, args ...interface{}) {
	DefaultLogger.Errorf(message, args...)
}

func Fatalf(message string, args ...interface{}) {
	DefaultLogger.Fatalf(message, args...)
}

type NativeLogger struct {
	stdoutLogger *log.Logger
	stderrLogger *log.Logger
}

func NewNativeLogger(out, err io.Writer) *NativeLogger {
	return &NativeLogger{
		stdoutLogger: log.New(out, "", log.Flags()),
		stderrLogger: log.New(err, "", log.Flags()),
	}
}

func (l *NativeLogger) Debugf(message string, args ...interface{}) {
	l.printLevel("DEBUG", message, args...)
}

func (l *NativeLogger) Infof(message string, args ...interface{}) {
	l.printLevel("INFO", message, args...)
}

func (l *NativeLogger) Warnf(message string, args ...interface{}) {
	l.printLevel("WARN", message, args...)
}

func (l *NativeLogger) Errorf(message string, args ...interface{}) {
	l.printLevel("ERROR", message, args...)
}

func (l *NativeLogger) Fatalf(message string, args ...interface{}) {
	l.printLevel("FATAL", message, args...)
	os.Exit(1)
}

func (l *NativeLogger) printLevel(level, message string, args ...interface{}) {
	line := fmt.Sprintf("[%s] %s", level, message)
	var log *log.Logger
	switch {
	case level == "ERROR" || level == "FATAL":
		log = l.stderrLogger
	default:
		log = l.stdoutLogger
	}
	log.Printf(line, args...)
}

package logging

import (
	"fmt"
	"io"
	"os"
)

type Logger interface {
	Debugf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Fatalf(message string, args ...interface{})
}

var DefaultLogger Logger = new(NativeLogger)

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
}

func (l NativeLogger) Debugf(message string, args ...interface{}) {
	l.printLevel(os.Stdout, "DEBUG", message, args...)
}

func (l NativeLogger) Infof(message string, args ...interface{}) {
	l.printLevel(os.Stdout, "INFO", message, args...)
}

func (l NativeLogger) Warnf(message string, args ...interface{}) {
	l.printLevel(os.Stdout, "WARN", message, args...)
}

func (l NativeLogger) Errorf(message string, args ...interface{}) {
	l.printLevel(os.Stderr, "ERROR", message, args...)
}

func (l NativeLogger) Fatalf(message string, args ...interface{}) {
	l.printLevel(os.Stderr, "FATAL", message, args...)
	os.Exit(1)
}

func (l NativeLogger) printLevel(out io.Writer, level, message string, args ...interface{}) {
	line := fmt.Sprintf("[%s] %s\n", level, message)
	fmt.Fprintf(out, line, args...)
}

type MutedLogger struct {
}

func (l MutedLogger) Debugf(message string, args ...interface{}) {
}

func (l MutedLogger) Infof(message string, args ...interface{}) {
}

func (l MutedLogger) Warnf(message string, args ...interface{}) {
}

func (l MutedLogger) Errorf(message string, args ...interface{}) {
}

func (l MutedLogger) Fatalf(message string, args ...interface{}) {
	os.Exit(1)
}

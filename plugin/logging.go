package plugin

import (
	"fmt"
	"os"
)

type Logger interface {
	Info(string)
	Error(string)
	Close() error
}

func NewLogger(name string) (Logger, error) {
	filename := fmt.Sprintf("%s.log", name)
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &logger{
		file: file,
	}, nil
}

type logger struct {
	file *os.File
}

func (l *logger) Info(value string) {
	fmt.Fprintln(l.file, fmt.Sprintf("INFO: %s\n", value))
}

func (l *logger) Error(value string) {
	fmt.Fprintln(l.file, fmt.Sprintf("ERROR: %s\n", value))
}

func (l *logger) Close() error {
	return l.file.Close()
}

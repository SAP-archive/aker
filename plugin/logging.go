package plugin

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func newLogWriter(name string, sink io.Writer) *logWriter {
	return &logWriter{
		name:   name,
		sink:   sink,
		buffer: new(bytes.Buffer),
	}
}

type logWriter struct {
	name   string
	sink   io.Writer
	buffer *bytes.Buffer
}

func (w *logWriter) Write(data []byte) (int, error) {
	w.buffer.Write(data)
	scanner := bufio.NewScanner(w.buffer)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return len(data), err
		}
		_, err := fmt.Fprintf(w.sink, "[%s]: %s\n", w.name, scanner.Text())
		if err != nil {
			return len(data), err
		}
	}
	return len(data), nil
}

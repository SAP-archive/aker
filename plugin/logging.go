package plugin

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func newLogWriter(name string, out io.Writer) io.Writer {
	return &logWriter{
		name:     name,
		delegate: out,
		buffer:   new(bytes.Buffer),
	}
}

type logWriter struct {
	name     string
	delegate io.Writer
	buffer   *bytes.Buffer
}

func (p *logWriter) Write(data []byte) (int, error) {
	p.buffer.Write(data)
	scanner := bufio.NewScanner(p.buffer)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return len(data), err
		}
		_, err := fmt.Fprintf(p.delegate, "[%s]: %s\n", p.name, scanner.Text())
		if err != nil {
			return len(data), err
		}
	}
	return len(data), nil
}

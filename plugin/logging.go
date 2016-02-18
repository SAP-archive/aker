package plugin

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

func newPluginLogWriter(name, level string, out io.Writer) io.Writer {
	return &pluginLogWriter{
		name:     name,
		level:    level,
		delegate: out,
		buffer:   new(bytes.Buffer),
	}
}

type pluginLogWriter struct {
	name     string
	level    string
	delegate io.Writer
	buffer   *bytes.Buffer
}

func (p *pluginLogWriter) Write(data []byte) (int, error) {
	p.buffer.Write(data)
	scanner := bufio.NewScanner(p.buffer)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return len(data), err
		}
		_, err := fmt.Fprintf(p.delegate, "[%s] %s: %s\n", p.name, p.level, scanner.Text())
		if err != nil {
			return len(data), err
		}
	}
	return len(data), nil
}

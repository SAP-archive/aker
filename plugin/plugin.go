package plugin

import (
	"errors"
	"net/http"
	"plugin"
)

// Plugin represents an Aker plugin.
type Plugin struct {
	handler http.Handler
	next    *Plugin
}

func (p *Plugin) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	rt := &responseTracker{w, false}
	p.handler.ServeHTTP(rt, req)
	if rt.done {
		return
	}
	p.next.ServeHTTP(w, req)
}

// InitFunc should be implemented by all plugins, so it is possible to initialize
// the http.Handlers that they provide.
type InitFunc func([]byte) (http.Handler, error)

// Open opens the plugin located at the given path.
func Open(path string, config []byte, next *Plugin) (*Plugin, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	sym, err := p.Lookup("Init")
	if err != nil {
		return nil, err
	}
	if init, ok := sym.(InitFunc); ok {
		h, err := init(config)
		if err != nil {
			return nil, err
		}
		return &Plugin{h, next}, nil
	}
	return nil, errors.New("plugin: missing Initer symbol")
}

type responseTracker struct {
	http.ResponseWriter
	done bool
}

func (w *responseTracker) Write(data []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(data)
}

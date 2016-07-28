package plugin

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"

	"github.infra.hana.ondemand.com/I061150/aker/socket"
)

// Open opens the specified plugin using the DefaultOpener.
func Open(name string, config []byte, next *Plugin) (*Plugin, error) {
	return DefaultOpener.Open(name, config, next)
}

// Opener wraps the basic Open plugin method.
type Opener interface {
	// Open should connect to and configure the specified plugin.
	Open(name string, config []byte, next *Plugin) (*Plugin, error)
}

// DefaultOpener is the default implementation of Opener.
var DefaultOpener = opener{}

type opener struct {
}

func (o opener) Open(name string, config []byte, next *Plugin) (*Plugin, error) {
	socketPath, err := socket.GetUniquePath("aker-plugin")
	if err != nil {
		return nil, err
	}

	setup, err := json.Marshal(&setup{
		SocketPath:        socketPath,
		ForwardSocketPath: next.SocketPath(),
		Configuration:     config,
	})
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(name)
	cmd.Stdin = bytes.NewReader(setup)
	cmd.Stdout = newLogWriter(name, os.Stdout)
	cmd.Stderr = newLogWriter(name, os.Stderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Plugin{
		socketPath: socketPath,
		Handler:    socket.ProxyHTTP(socketPath),
	}, nil
}

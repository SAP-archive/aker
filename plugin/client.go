package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	"github.infra.hana.ondemand.com/I061150/aker/socket"
)

//go:generate counterfeiter . Opener

// Plugin represents an Aker plugin.
type Plugin struct {
	http.Handler
	socketPath string
}

// SocketPath returns the path of the socket that the plugin is binded to.
func (p *Plugin) SocketPath() string {
	if p == nil {
		return ""
	}
	return p.socketPath
}

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

type opener struct{}

func (o opener) Open(name string, config []byte, next *Plugin) (*Plugin, error) {
	socketPath, err := socket.GetUniqueSocketPath("aker-plugin")
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
		Handler:    socket.Proxy(socketPath),
	}, nil
}

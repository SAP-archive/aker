package plugin

import (
	"net/http"
	"os"
)

// Plugin represents an Aker plugin.
type Plugin struct {
	http.Handler
	socketPath string
	process    *os.Process
}

// SocketPath returns the path of the socket that the plugin is binded to.
func (p *Plugin) SocketPath() string {
	if p == nil {
		return ""
	}
	return p.socketPath
}

// Close releases all resources allocated by the plugin.
func (p *Plugin) Close() error {
	if err := p.process.Signal(os.Interrupt); err != nil {
		return err
	}
	_, err := p.process.Wait()
	return err
}

type setup struct {
	SocketPath        string `json:"socket_path"`
	ForwardSocketPath string `json:"forward_socket_path"`
	Configuration     []byte `json:"configuration"`
}

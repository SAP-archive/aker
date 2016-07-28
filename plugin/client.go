package plugin

import "net/http"

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

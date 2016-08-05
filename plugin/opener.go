package plugin

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.infra.hana.ondemand.com/cloudfoundry/aker/socket"
)

// Open opens the specified plugin using the DefaultOpener.
func Open(name string, config []byte, next *Plugin) (*Plugin, error) {
	return DefaultOpener.Open(name, config, next)
}

// DefaultOpener redirects the plugin's stdout and stderr to the calling process's
// stdout and stderr.
var DefaultOpener = &Opener{
	PluginStdout: os.Stdout,
	PluginStderr: os.Stderr,
}

// Opener opens a plugin as a subprocess.
type Opener struct {
	// Stdout of the child process is redirected to PluginStdout.
	PluginStdout io.Writer
	// Stderr of the child process is redirected to PluginStderr.
	PluginStderr io.Writer
}

func (o *Opener) Open(name string, config []byte, next *Plugin) (*Plugin, error) {
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
	cmd.Stdout = newLogWriter(name, o.PluginStdout)
	cmd.Stderr = newLogWriter(name, o.PluginStderr)
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Plugin{
		socketPath: socketPath,
		Handler:    socket.ProxyHTTP(socketPath),
		process:    cmd.Process,
	}, nil
}

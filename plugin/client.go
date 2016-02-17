package plugin

import (
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

const connectAttempts = 5
const connectAttemptInterval = time.Second

// Start starts the plugin with the specified
// name as a new process and prepares it for listening
func Start(name string) error {
	cmd := exec.Command(name)
	cmd.Stdin = strings.NewReader(sockPath(name))
	cmd.Stdout = newPluginLogWriter(name, "INFO", os.Stdout)
	cmd.Stderr = newPluginLogWriter(name, "ERROR", os.Stderr)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

// Connect tries to open a connection to the plugin with
// the specified name.
// It is important that the plugin was started with Start
// beforehand.
func Connect(name string) (api.Plugin, error) {
	var err error
	var conn net.Conn
	for attempts := 1; attempts <= connectAttempts; attempts++ {
		conn, err = net.Dial("unix", sockPath(name))
		if err == nil {
			const maxInt = int(^uint(0) >> 1)
			peer := connection.NewPeer(conn, conn, maxInt/2)
			return connection.NewPluginClient(peer), nil
		}
		if attempts < connectAttempts {
			time.Sleep(connectAttemptInterval)
		}
	}
	return nil, err
}

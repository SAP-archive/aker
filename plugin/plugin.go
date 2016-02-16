package plugin

import (
	"net"
	"os"
	"os/exec"
	"time"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

func ListenAndServe(plugin api.Plugin) {
	listener, err := net.Listen("unix", "/tmp/aker-plugin.sock")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			peer := connection.NewPeer(conn, conn, 1)
			connection.ServePlugin(peer, plugin)
		}(conn)
	}
}

// FIXME: Only one instance of plugin supported at the moment
func Open(name string) (api.Plugin, error) {
	cmd := exec.Command(name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second) // FIXME

	conn, err := net.Dial("unix", "/tmp/aker-plugin.sock")
	if err != nil {
		return nil, err
	}

	const maxInt = int(^uint(0) >> 1)
	peer := connection.NewPeer(conn, conn, maxInt/2)
	return connection.NewPluginClient(peer), nil
}

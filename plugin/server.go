package plugin

import (
	"net"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

type PluginFactory func() (api.Plugin, error)

func ListenAndServe(name string, factory PluginFactory) error {
	listener, err := net.Listen("unix", sockPath(name))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		plugin, err := factory()
		if err != nil {
			return err
		}

		go func(conn net.Conn, plug api.Plugin) {
			peer := connection.NewPeer(conn, conn, 1)
			connection.ServePlugin(peer, plug)
		}(conn, plugin)
	}
}

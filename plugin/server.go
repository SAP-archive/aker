package plugin

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

type PluginFactory func() (api.Plugin, error)

func ListenAndServe(factory PluginFactory) error {
	socketPath, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	fmt.Printf("Socket path: %s\n", string(socketPath))

	listener, err := net.Listen("unix", string(socketPath))
	if err != nil {
		return err
	}
	defer listener.Close()

	go handlePluginConnections(listener, factory)

	osChannel := make(chan os.Signal)
	signal.Notify(osChannel, syscall.SIGTERM, syscall.SIGINT)
	for range osChannel {
		return nil
	}
	return nil
}

func handlePluginConnections(listener net.Listener, factory PluginFactory) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		plugin, err := factory()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn, plug api.Plugin) {
			peer := connection.NewPeer(conn, conn, 1)
			connection.ServePlugin(peer, plug)
		}(conn, plugin)
	}
}

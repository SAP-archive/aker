package plugin

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

func ListenAndServe(plugin api.Plugin) {
	peer := connection.NewPeer(os.Stdin, os.Stdout, 1)
	connection.ServePlugin(peer, plugin)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	for range ch {
		os.Exit(0)
	}
}

func Open(name string) (api.Plugin, error) {
	cmd := exec.Command(name)
	pluginIn, cmdOut := io.Pipe()
	cmd.Stdout = cmdOut
	cmdIn, pluginOut := io.Pipe()
	cmd.Stdin = cmdIn
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	const maxInt = int(^uint(0) >> 1)
	peer := connection.NewPeer(pluginIn, pluginOut, maxInt/2)
	return connection.NewPluginClient(peer), nil
}

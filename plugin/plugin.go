package plugin

import (
	"io"
	"os"
	"os/exec"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/connection"
)

func ListenAndServe(plugin api.Plugin) {
	peer := connection.NewPeer(os.Stdin, os.Stdout, 10)
	connection.ServePlugin(peer, plugin)
}

func Open(name string) (api.Plugin, error) {
	cmd := exec.Command(name)
	pluginIn, cmdOut := io.Pipe()
	cmd.Stdout = cmdOut
	cmdIn, pluginOut := io.Pipe()
	cmd.Stdin = cmdIn
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	peer := connection.NewPeer(pluginIn, pluginOut, 100000000)
	return connection.NewPluginClient(peer), nil
}

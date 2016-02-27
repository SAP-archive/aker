package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"

	"github.wdf.sap.corp/I061150/aker/config"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	cfg, err := config.LoadFromFile("config.yaml")
	if err != nil {
		panic(err)
	}

	for _, endpoint := range cfg.Endpoints {
		leadingPlugin, err := buildPluginChain(endpoint.Plugins)
		if err != nil {
			panic(err)
		}
		http.Handle(endpoint.Path, leadingPlugin)
	}

	fmt.Println("Starting HTTP listener...")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err = http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("HTTP Listener failed with '%s'.\n", err)
		os.Exit(1)
	}
}

func buildPluginChain(references []config.PluginReferenceConfig) (http.Handler, error) {
	index := len(references) - 1
	lastPlugin, err := buildPlugin(references[index], nil)
	if err != nil {
		return nil, err
	}
	for index > 0 {
		index--
		lastPlugin, err = buildPlugin(references[index], lastPlugin)
		if err != nil {
			return nil, err
		}
	}
	return lastPlugin, nil
}

func buildPlugin(cfg config.PluginReferenceConfig, next *plugin.Plugin) (*plugin.Plugin, error) {
	cfgData, err := candiedyaml.Marshal(cfg.Config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Opening plugin: %s\n", cfg.Name)
	plug, err := plugin.Open(cfg.Name, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

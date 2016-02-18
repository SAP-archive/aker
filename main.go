package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.wdf.sap.corp/I061150/aker/config"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	cfg, err := config.LoadFromFile("config.json")
	if err != nil {
		panic(err)
	}

	for _, handler := range cfg.Handlers {
		firstFilter, err := buildPluginChain(handler.Filters)
		if err != nil {
			panic(err)
		}
		http.Handle(handler.Path, firstFilter)
	}

	fmt.Println("Starting HTTP listener...")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err = http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("HTTP Listener failed with '%s'.\n", err)
		os.Exit(1)
	}
}

func buildPluginChain(filterConfigs []config.FilterConfig) (http.Handler, error) {
	index := len(filterConfigs) - 1
	lastPlugin, err := buildPlugin(filterConfigs[index], nil)
	if err != nil {
		return nil, err
	}
	for index > 0 {
		index--
		lastPlugin, err = buildPlugin(filterConfigs[index], lastPlugin)
		if err != nil {
			return nil, err
		}
	}
	return lastPlugin, nil
}

func buildPlugin(cfg config.FilterConfig, next *plugin.Plugin) (*plugin.Plugin, error) {
	cfgData, err := json.Marshal(cfg.PluginConfig)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Opening plugin: %s\n", cfg.PluginName)
	plug, err := plugin.Open(cfg.PluginName, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

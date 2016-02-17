package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.wdf.sap.corp/I061150/aker/api"
	"github.wdf.sap.corp/I061150/aker/config"
	"github.wdf.sap.corp/I061150/aker/flow"
	"github.wdf.sap.corp/I061150/aker/plugin"
)

func main() {
	cfg, err := config.LoadFromFile("config.json")
	if err != nil {
		panic(err)
	}

	pluginNames := uniquePluginNames(cfg)
	for _, pluginName := range pluginNames {
		fmt.Printf("Starting plugin '%s'...\n", pluginName)
		if err := plugin.Start(pluginName); err != nil {
			fmt.Printf("Failed to start plugin '%s' due to '%s'!\n", pluginName, err)
			os.Exit(1)
		}
		fmt.Printf("Plugin '%s' started successfully.\n", pluginName)
	}

	for _, handler := range cfg.Handlers {
		filters, err := buildFilters(handler.Filters)
		if err != nil {
			panic(err)
		}
		http.Handle(handler.Path, flow.NewChainedFilterHandler(filters...))
	}

	fmt.Println("Starting HTTP listener...")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err = http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("HTTP Listener failed with '%s'.\n", err)
		os.Exit(1)
	}
}

func uniquePluginNames(cfg config.Config) []string {
	names := map[string]struct{}{}
	for _, handler := range cfg.Handlers {
		for _, filter := range handler.Filters {
			names[filter.PluginName] = struct{}{}
		}
	}
	index := 0
	result := make([]string, len(names))
	for name := range names {
		result[index] = name
		index++
	}
	return result
}

func buildFilters(filterConfigs []config.FilterConfig) ([]api.Plugin, error) {
	result := []api.Plugin{}
	for _, fltrConfig := range filterConfigs {
		filter, err := buildFilterWithConfig(fltrConfig)
		if err != nil {
			return nil, err
		}
		result = append(result, filter)
	}
	return result, nil
}

func buildFilterWithConfig(cfg config.FilterConfig) (api.Plugin, error) {
	fmt.Printf("Connecting to plugin: %s\n", cfg.PluginName)
	plug, err := plugin.Connect(cfg.PluginName)
	if err != nil {
		return nil, err
	}
	cfgData, err := json.Marshal(cfg.PluginConfig)
	if err != nil {
		return nil, err
	}
	if err := plug.Config(cfgData); err != nil {
		return nil, err
	}
	return plug, nil
}

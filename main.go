package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	}
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
	fmt.Printf("Opening plugin %s\n", cfg.PluginName)
	plug, err := plugin.Open(cfg.PluginName)
	if err != nil {
		return nil, err
	}
	fmt.Println("Configuring plugin...")
	cfgData, err := json.Marshal(cfg.PluginConfig)
	if err != nil {
		return nil, err
	}
	if err := plug.Config(cfgData); err != nil {
		return nil, err
	}
	fmt.Println("Done!")
	return plug, nil
}

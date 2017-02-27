package endpoint

import (
	"net/http"
	"os"

	"github.com/SAP/aker/config"
	"github.com/SAP/aker/logging"
	"github.com/SAP/aker/plugin"
	"github.com/SAP/gologger"
)

// Handler represents Aker endpoint.
type Handler struct {
	http.Handler
	path string
}

// NewHandler creates new endpoint handler. It opens all plugins specified
// by endpoint.
func NewHandler(endpoint config.Endpoint) (*Handler, error) {
	if endpoint.Path == "" {
		return nil, InvalidPathError("")
	}
	if endpoint.Plugins == nil || len(endpoint.Plugins) == 0 {
		return nil, NoPluginsErr
	}

	pluginChain, err := buildChain(endpoint.Plugins)
	if err != nil {
		return nil, err
	}

	if endpoint.Audit {
		pluginChain = logging.Handler(os.Stdout, pluginChain)
	}

	return &Handler{
		Handler: pluginChain,
		path:    endpoint.Path,
	}, nil
}

func buildChain(references []config.PluginReference) (http.Handler, error) {
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

func buildPlugin(reference config.PluginReference, next *plugin.Plugin) (*plugin.Plugin, error) {
	cfgData, err := plugin.MarshalConfig(reference.Config)
	if err != nil {
		return nil, err
	}

	gologger.Infof("Opening plugin: %q", reference.Name)
	plug, err := plugin.Open(reference.Name, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

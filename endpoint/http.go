package endpoint

import (
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"

	"github.infra.hana.ondemand.com/I061150/aker/config"
	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
)

type Handler struct {
	path        string
	pluginChain http.Handler
}

func NewHandler(endpoint config.Endpoint) (*Handler, error) {
	pluginChain, err := buildPluginChain(endpoint.Plugins)
	if err != nil {
		return nil, err
	}

	if endpoint.Audit {
		pluginChain = logging.Handler(os.Stdout, pluginChain)
	}

	return &Handler{
		path:        endpoint.Path,
		pluginChain: pluginChain,
	}, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.pluginChain.ServeHTTP(w, req)
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

	logging.Infof("Opening plugin: %q", cfg.Name)
	plug, err := plugin.Open(cfg.Name, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

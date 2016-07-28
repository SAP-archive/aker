package endpoint

import (
	"net/http"
	"os"

	"github.infra.hana.ondemand.com/I061150/aker/config"
	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
)

// Handler represents Aker endpoint.
type Handler struct {
	path        string
	plugin      plugin.Opener
	pluginChain http.Handler
}

// NewHandler creates new endpoint handler. It opens all plugins specified
// by endpoint using the provided Opener.
func NewHandler(endpoint config.Endpoint, opener plugin.Opener) (*Handler, error) {
	if endpoint.Path == "" {
		return nil, InvalidPathError("")
	}
	if endpoint.Plugins == nil || len(endpoint.Plugins) == 0 {
		return nil, NoPluginsErr
	}

	chainBuilder := chainBuilder{opener}
	pluginChain, err := chainBuilder.build(endpoint.Plugins)
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

// ServeHTTP routes the incoming http.Request through the chain of aker plugins.
func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.pluginChain.ServeHTTP(w, req)
}

type chainBuilder struct {
	plugin plugin.Opener
}

func (b *chainBuilder) build(references []config.PluginReference) (http.Handler, error) {
	index := len(references) - 1
	lastPlugin, err := b.buildPlugin(references[index], nil)
	if err != nil {
		return nil, err
	}
	for index > 0 {
		index--
		lastPlugin, err = b.buildPlugin(references[index], lastPlugin)
		if err != nil {
			return nil, err
		}
	}
	return lastPlugin, nil
}

func (b *chainBuilder) buildPlugin(reference config.PluginReference, next *plugin.Plugin) (*plugin.Plugin, error) {
	cfgData, err := plugin.MarshalConfig(reference.Config)
	if err != nil {
		return nil, err
	}

	logging.Infof("Opening plugin: %q", reference.Name)
	plug, err := b.plugin.Open(reference.Name, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

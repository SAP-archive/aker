package endpoint

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"

	"github.infra.hana.ondemand.com/I061150/aker/config"
	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
)

type InvalidPathError string

func (e InvalidPathError) Error() string {
	return fmt.Sprintf("invalid endpoint path: %q", e)
}

var NoPluginsErr = errors.New("no plugins specified")

type Handler struct {
	path        string
	plugin      plugin.Opener
	pluginChain http.Handler
}

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

func (b *chainBuilder) buildPlugin(cfg config.PluginReference, next *plugin.Plugin) (*plugin.Plugin, error) {
	cfgData, err := candiedyaml.Marshal(cfg.Config)
	if err != nil {
		return nil, err
	}

	logging.Infof("Opening plugin: %q", cfg.Name)
	plug, err := b.plugin.Open(cfg.Name, cfgData, next)
	if err != nil {
		return nil, err
	}
	return plug, nil
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"

	"github.infra.hana.ondemand.com/I061150/aker/config"
	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"
	"github.infra.hana.ondemand.com/I061150/aker/uuid"
)

var configLocationFlag = flag.String(
	"config",
	"config.yml",
	"Specifies the configuration file location. By default this is './config.yml'.",
)

func main() {
	flag.Parse()

	cfg, err := config.LoadFromFile(*configLocationFlag)
	if err != nil {
		logging.Fatalf("Failed to load configuration due to %q", err.Error())
	}
	mux := http.NewServeMux()
	for _, endpoint := range cfg.Endpoints {
		leadingPlugin, err := buildPluginChain(endpoint.Plugins)
		if err != nil {
			logging.Fatalf("Failed to build plugin chain due to %q", err.Error())
		}
		if endpoint.Audit {
			leadingPlugin = logging.Handler(os.Stdout, leadingPlugin)
		}
		mux.Handle(endpoint.Path, leadingPlugin)
	}

	logging.Infof("Starting HTTP listener...")
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	handler := NewHeaderSticker(mux, map[string]func() string{
		"X-Aker-Request-Id": func() string {
			uid, _ := uuid.Random()
			return uid.String()
		},
	})

	if err = http.ListenAndServe(addr, handler); err != nil {
		logging.Fatalf("HTTP Listener failed with %q", err.Error())
	}
}

type HeaderSticker struct {
	http.Handler
	headers map[string]func() string
}

func NewHeaderSticker(h http.Handler, headers map[string]func() string) *HeaderSticker {
	return &HeaderSticker{
		Handler: h,
		headers: headers,
	}
}

func (s *HeaderSticker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for header, valueFunc := range s.headers {
		req.Header.Add(header, valueFunc())
	}
	s.Handler.ServeHTTP(w, req)
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

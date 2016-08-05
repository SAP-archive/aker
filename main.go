package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.infra.hana.ondemand.com/cloudfoundry/aker/config"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/endpoint"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/logging"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/uuid"
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
	for _, endpointCfg := range cfg.Endpoints {
		endpointHandler, err := endpoint.NewHandler(endpointCfg, plugin.DefaultOpener)
		if err != nil {
			logging.Fatalf("Failed to build plugin chain due to %q", err.Error())
		}
		mux.Handle(endpointCfg.Path, endpointHandler)
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
		value := valueFunc()
		req.Header.Add(header, value)
		w.Header().Add(header, value)
	}
	s.Handler.ServeHTTP(w, req)
}

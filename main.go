package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/SAP/aker/config"
	"github.com/SAP/aker/endpoint"
	"github.com/SAP/aker/uuid"
	"github.com/SAP/gologger"
)

var configLocationFlag = flag.String(
	"config",
	"config.yml",
	"Specifies the configuration file location.",
)

func main() {
	flag.Parse()

	cfg, err := config.LoadFromFile(*configLocationFlag)
	if err != nil {
		gologger.Fatalf("Failed to load configuration due to %q", err.Error())
	}
	mux := http.NewServeMux()
	for _, endpointCfg := range cfg.Endpoints {
		endpointHandler, err := endpoint.NewHandler(endpointCfg)
		if err != nil {
			gologger.Fatalf("Failed to build plugin chain due to %q", err.Error())
		}
		mux.Handle(endpointCfg.Path, endpointHandler)
	}

	handler := NewHeaderSticker(mux, map[string]func() string{
		"X-Aker-Request-Id": func() string {
			uid, _ := uuid.Random()
			return uid.String()
		},
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	gologger.Infof("Starting HTTP listener...")
	if err = srv.ListenAndServe(); err != nil {
		gologger.Fatalf("HTTP Listener failed with %q", err.Error())
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

package plugin

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.infra.hana.ondemand.com/cloudfoundry/aker/socket"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"
)

//go:generate counterfeiter . Socket
//go:generate counterfeiter . HTTPServer
//go:generate counterfeiter . Notifier

type HandlerFactory func(config []byte) (http.Handler, error)

// HTTPServer represents a HTTP server that could be started and then stopped.
type HTTPServer interface {
	Start() error
	Stop() error
}

// Socket enables running HTTP using unix domain socket transport.
type Socket interface {
	// ProxyHTTP should return a Handler that proxies all request to the
	// provided unix domain socket path.
	ProxyHTTP(socketPath string) http.Handler
	// NewHTTPServer creates a HTTPServer on the provided unix socket path.
	// The server uses the passed handler for serving HTTP requests.
	// The server is not started. It is caller's responsibility to start it.
	NewHTTPServer(path string, h http.Handler) HTTPServer
}

// Notifier notifies on OS signals.
type Notifier interface {
	// Notify relays incoming signals to c. If no
	// signals are provided, all incoming signals will be relayed to c.
	// Otherwise, just the provided signals will.
	Notify(c chan<- os.Signal, sig ...os.Signal)
}

// Server provides a functionality to run a Plugin.
type Server struct {
	config io.Reader
	socket Socket
	signal Notifier
	log    gologger.Logger
}

// NewServer returns a brand new server.
func NewServer(config io.Reader, log gologger.Logger, socket Socket, signal Notifier) *Server {
	return &Server{
		config: config,
		socket: socket,
		signal: signal,
		log:    log,
	}
}

// socketProxy is a proxy for the socket package.
type socketProxy struct{}

// ProxyHTTP calls ProxyHTTP method from socket package.
func (s socketProxy) ProxyHTTP(socketPath string) http.Handler {
	return socket.ProxyHTTP(socketPath)
}

func (s socketProxy) NewHTTPServer(path string, handler http.Handler) HTTPServer {
	return socket.NewHTTPServer(path, handler)
}

// notifier calls os.Notify
type notifier struct{}

func (n notifier) Notify(c chan<- os.Signal, sig ...os.Signal) {
	signal.Notify(c, sig...)
}

// DefaultServer is a server with default configuration.
// One should not need to create a different server. If you're creating new
// server and reading this, you're doing something wrong.
//
// It uses os.Stdin for reading configuration.
// It uses the gologger.DefaultLogger to log.
// It uses the socket package for HTTP over unix domain sockets.
// It uses the signal package from the standard library for signal handling.
var DefaultServer = NewServer(os.Stdin, gologger.DefaultLogger, socketProxy{}, notifier{})

// ListenAndServeHTTP starts the http.Handler returned by the factory as plugin.
func (s *Server) ListenAndServeHTTP(factory HandlerFactory) error {
	var setup setup
	decoder := json.NewDecoder(s.config)
	if err := decoder.Decode(&setup); err != nil {
		return &ConfigDecodeError{err}
	}

	handler, err := factory(setup.Configuration)
	if err != nil {
		return err
	}
	if setup.ForwardSocketPath != "" {
		handler = &forwardHandler{
			current: handler,
			next:    s.socket.ProxyHTTP(setup.ForwardSocketPath),
		}
	}

	s.log.Infof("Listening on socket: %s\n", setup.SocketPath)

	server := s.socket.NewHTTPServer(setup.SocketPath, handler)
	if err := server.Start(); err != nil {
		s.log.Fatalf("Error starting server: %v\n", err)
	}
	defer server.Stop()

	c := make(chan os.Signal)
	s.signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	s.log.Infof("Exiting due to: %v\n", sig)
	return nil
}

// ListenAndServeHTTP calls ListenAndServeHTTP of the DefaultServer.
func ListenAndServeHTTP(factory HandlerFactory) error {
	return DefaultServer.ListenAndServeHTTP(factory)
}

type responseTracker struct {
	http.ResponseWriter
	done bool
}

func (w *responseTracker) Write(data []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(data)
}

func (w *responseTracker) WriteHeader(status int) {
	w.done = true
	w.ResponseWriter.WriteHeader(status)
}

type forwardHandler struct {
	current http.Handler
	next    http.Handler
}

func (h *forwardHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	respTracker := &responseTracker{
		ResponseWriter: resp,
		done:           false,
	}
	h.current.ServeHTTP(respTracker, req)
	if respTracker.done {
		return
	}
	h.next.ServeHTTP(resp, req)
}

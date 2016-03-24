package plugin

import (
	"encoding/json"
	"net/http"
	"os"

	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/socket"
)

type HandlerFactory func(config []byte) (http.Handler, error)

func ListenAndServe(factory HandlerFactory) error {
	var setup setup
	decoder := json.NewDecoder(os.Stdin)
	if err := decoder.Decode(&setup); err != nil {
		return err
	}

	handler, err := factory(setup.Configuration)
	if err != nil {
		return err
	}
	if setup.ForwardSocketPath != "" {
		handler = newForwardHandler(handler, setup.ForwardSocketPath)
	}

	logging.Infof("Listening on socket: %s\n", setup.SocketPath)
	return socket.ListenAndServe(setup.SocketPath, handler)
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

func newForwardHandler(current http.Handler, nextSocket string) http.Handler {
	return &forwardHandler{
		current: current,
		next:    socket.Proxy(nextSocket),
	}
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

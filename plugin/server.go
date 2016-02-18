package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.wdf.sap.corp/I061150/aker/socket"
)

type HandlerFactory func(config []byte) (http.Handler, error)

func ListenAndServe(factory HandlerFactory) error {
	var setup pluginSetup
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

	fmt.Printf("Listening on socket: %s\n", setup.SocketPath)
	return socket.ListenAndServe(setup.SocketPath, handler)
}

type responseWrapper struct {
	http.ResponseWriter
	done bool
}

func (w *responseWrapper) Write(data []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(data)
}

func (w *responseWrapper) WriteHeader(status int) {
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
	respWrapper := &responseWrapper{
		ResponseWriter: resp,
		done:           false,
	}
	h.current.ServeHTTP(respWrapper, req)
	if respWrapper.done {
		return
	}
	h.next.ServeHTTP(resp, req)
}

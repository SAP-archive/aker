package flow

import (
	"net/http"
	"net/url"

	"github.wdf.sap.corp/I061150/aker/api"
)

type chainedFilterHandler struct {
	filters []api.Plugin
}

func NewChainedFilterHandler(filters ...api.Plugin) http.Handler {
	return &chainedFilterHandler{
		filters: filters,
	}
}

func (h *chainedFilterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	req := &requestWrapper{
		original: request,
	}
	resp := &responseWrapper{
		original: response,
	}
	data := &dataWrapper{
		stringValues: make(map[string]string),
	}
	context := api.Context{
		Request:  req,
		Response: resp,
		Data:     data,
	}
	for _, filter := range h.filters {
		done := filter.Process(context)
		if done {
			break
		}
	}
}

type requestWrapper struct {
	original *http.Request
}

func (w *requestWrapper) URL() *url.URL {
	return w.original.URL
}

func (w *requestWrapper) Method() string {
	return w.original.Method
}

func (w *requestWrapper) Host() string {
	return w.original.Host
}

func (w *requestWrapper) ContentLength() int {
	return int(w.original.ContentLength) // FIXME: Dangerous conversion
}

func (w *requestWrapper) Headers() map[string][]string {
	return w.original.Header
}

func (w *requestWrapper) Header(name string) string {
	return w.original.Header.Get(name)
}

func (w *requestWrapper) Read(data []byte) (int, error) {
	return w.original.Body.Read(data)
}

func (w *requestWrapper) Close() error {
	return w.original.Body.Close()
}

type responseWrapper struct {
	original http.ResponseWriter
}

func (w *responseWrapper) SetHeader(name string, values []string) {
	w.original.Header().Del(name)
	for _, value := range values {
		w.original.Header().Add(name, value)
	}
}

func (w *responseWrapper) WriteStatus(status int) {
	w.original.WriteHeader(status)
}

func (w *responseWrapper) Write(data []byte) (int, error) {
	return w.original.Write(data)
}

type dataWrapper struct {
	stringValues map[string]string
}

func (w *dataWrapper) SetString(name, value string) {
	w.stringValues[name] = value
}

func (w *dataWrapper) String(name string) string {
	return w.stringValues[name]
}

package adapter

import (
	"net/http"

	"github.wdf.sap.corp/I061150/aker/api"
)

type ResponseWriterAdapter interface {
	http.ResponseWriter
	Flush()
}

func NewResponseWriterAdapter(delegate api.Response) ResponseWriterAdapter {
	return &responseWrapper{
		delegate:       delegate,
		headers:        make(map[string][]string),
		headersWritten: false,
	}
}

type responseWrapper struct {
	delegate       api.Response
	headers        http.Header
	headersWritten bool
}

func (w *responseWrapper) Header() http.Header {
	return w.headers
}

func (w *responseWrapper) WriteHeader(status int) {
	for name, value := range w.headers {
		w.delegate.SetHeader(name, value)
	}
	w.delegate.WriteStatus(status)
	w.headersWritten = true
}

func (w *responseWrapper) Write(data []byte) (int, error) {
	return w.delegate.Write(data)
}

func (w *responseWrapper) Flush() {
	if !w.headersWritten {
		w.WriteHeader(http.StatusOK)
	}
}

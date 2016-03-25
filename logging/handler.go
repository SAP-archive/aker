package logging

import (
	"io"
	"net/http"
	"time"
)

//go:generate counterfeiter . ResponseRecorder

// DefaultFormatter is the default implementation of Formatter and is used by
// DefaultLoggingHandlerFactory.
var DefaultFormatter = defaultFormatter{}

// TimeProviderFunc represents a factory function that returns the current time
type TimeProviderFunc func() time.Time

// HandlerFactory constructs handlers that logs each request.
//
// The format of each log entry is determined via the Formatter.
//
// The StartedAt and FinishedAt fields of each AccessEntry are determined via
// the specified time provider. This is useful at least for testing purposes.
type HandlerFactory struct {
	TimeProvider TimeProviderFunc
	Formatter    Formatter
}

// LoggingHandler returns new http.Handler that logs.
func (f *HandlerFactory) LoggingHandler(log io.Writer, h http.Handler) http.Handler {
	return &loggingHandler{
		Handler: h,
		log:     log,
		now:     f.TimeProvider,
		format:  f.Formatter,
	}
}

// DefaultLoggingHandlerFactory is the default implementation of
// HandlerFactory.
var DefaultLoggingHandlerFactory = HandlerFactory{
	TimeProvider: time.Now,
	Formatter:    DefaultFormatter,
}

// Handler returns http.Handler that logs using the DefaultFormatter.
func Handler(log io.Writer, h http.Handler) http.Handler {
	return DefaultLoggingHandlerFactory.LoggingHandler(log, h)
}

type loggingHandler struct {
	http.Handler
	format Formatter
	log    io.Writer
	now    TimeProviderFunc
}

// ServeHTTP calls the ServeHTTP of the underlying handler and logs.
func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	wrec := &responseRecorder{ResponseWriter: w}
	defer func(startedAt time.Time) {
		e := &AccessEntry{
			StartedAt:  startedAt,
			FinishedAt: h.now(),
			Request:    req,
			Response:   wrec,
		}
		h.log.Write([]byte(h.format.Format(e)))
	}(h.now())

	h.Handler.ServeHTTP(wrec, req)
}

// AccessEntry contains information about a processed HTTP request.
// It should not be modified.
type AccessEntry struct {
	Request    *http.Request
	Response   ResponseRecorder
	StartedAt  time.Time
	FinishedAt time.Time
}

// ResponseRecorder represents http.ResponseWriter which keeps track of the
// response status code and response body size (in bytes).
type ResponseRecorder interface {
	http.ResponseWriter
	Status() int
	Size() int
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *responseRecorder) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.size += n
	return n, err
}

func (r *responseRecorder) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.status = status
}

// Status returns the HTTP status code of the response.
func (r *responseRecorder) Status() int {
	return r.status
}

// Size returns the length in bytes of the response body.
func (r *responseRecorder) Size() int {
	return r.size
}

package logging

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type TimeProviderFunc func() time.Time

type LoggingHandlerFactory struct {
	TimeProvider TimeProviderFunc
}

func (f *LoggingHandlerFactory) LoggingHandler(log io.Writer, h http.Handler) http.Handler {
	return &loggingHandler{
		Handler: h,
		log:     log,
		now:     f.TimeProvider,
	}
}

func LoggingHandler(log io.Writer, h http.Handler) http.Handler {
	f := &LoggingHandlerFactory{time.Now}
	return f.LoggingHandler(log, h)
}

type loggingHandler struct {
	http.Handler
	log io.Writer
	now TimeProviderFunc
}

func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	wrec := &responseRecorder{ResponseWriter: w}
	defer func(startedAt time.Time) {
		e := &accessEntry{
			startedAt:  startedAt,
			finishedAt: h.now(),
			req:        req,
			resp:       wrec,
		}
		e.WriteTo(h.log)
	}(h.now())

	h.Handler.ServeHTTP(wrec, req)
}

type accessEntry struct {
	req        *http.Request
	resp       *responseRecorder
	startedAt  time.Time
	finishedAt time.Time
}

func (e *accessEntry) String() string {
	return fmt.Sprintf(`%s - [%s] "%s %s %s" %d %d %d "%s" "%s" %s aker_request_id:%s response_time:%f`+"\n",
		e.req.URL.Host,
		e.startedAt.Format("02/01/2006:15:04:05 -0700"),
		e.req.Method,
		e.req.URL.RequestURI(),
		e.req.Proto,
		e.resp.status,
		e.req.ContentLength,
		e.resp.size,
		e.formatHeader("Referer"),
		e.formatHeader("User-Agent"),
		e.req.RemoteAddr,
		e.formatHeader("X-Aker-Request-Id"),
		e.finishedAt.Sub(e.startedAt).Seconds(),
	)
}

func (e *accessEntry) WriteTo(w io.Writer) (int64, error) {
	count, err := w.Write([]byte(e.String()))
	return int64(count), err
}

func (e *accessEntry) formatHeader(name string) string {
	v := e.req.Header.Get(name)
	if v == "" {
		return "-"
	}
	return v
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

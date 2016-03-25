package logging

import (
	"fmt"
	"net/http"
)

//go:generate counterfeiter . Formatter

// Formatter is an interface that wraps Format method for an access log entry.
// It should be safe for concurrent use.
type Formatter interface {
	// Format should return string representation of the passed AccessEntry.
	// It should not modify the AccessEntry.
	Format(*AccessEntry) string
}

type defaultFormatter struct{}

func (f defaultFormatter) Format(e *AccessEntry) string {
	return fmt.Sprintf(`%s - [%s] "%s %s %s" %d %d %d "%s" "%s" %s aker_request_id:%s response_time:%f`+"\n",
		e.Request.URL.Host,
		e.StartedAt.Format("02/01/2006:15:04:05 -0700"),
		e.Request.Method,
		e.Request.URL.RequestURI(),
		e.Request.Proto,
		e.Response.Status(),
		e.Request.ContentLength,
		e.Response.Size(),
		formatHeader(e.Request, "Referer"),
		formatHeader(e.Request, "User-Agent"),
		e.Request.RemoteAddr,
		formatHeader(e.Request, "X-Aker-Request-Id"),
		e.FinishedAt.Sub(e.StartedAt).Seconds())
}
func formatHeader(req *http.Request, name string) string {
	v := req.Header.Get(name)
	if v == "" {
		return "-"
	}
	return v
}

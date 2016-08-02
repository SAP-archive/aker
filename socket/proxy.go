package socket

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

const retryCount = 5
const retryInterval = time.Second

// ProxyHTTP proxies all requests to the specified socket path.
func ProxyHTTP(socketPath string) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http"
		},
		Transport: &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				var err error
				var conn net.Conn
				for i := 1; i <= retryCount; i++ {
					if conn, err = net.Dial("unix", socketPath); err == nil {
						return conn, nil
					}
					if i < retryCount {
						time.Sleep(retryInterval)
					}
				}
				return nil, err
			},
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

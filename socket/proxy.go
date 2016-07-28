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
			req.URL.Host = "localhost"
		},
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				var err error
				var con net.Conn
				for i := 1; i <= retryCount; i++ {
					if con, err = net.Dial("unix", socketPath); err == nil {
						return con, nil
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

package socket

import (
	"net"
	"net/http"
)

func ListenAndServe(path string, handler http.Handler) error {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return err
	}
	defer listener.Close()

	server := &http.Server{
		Handler: handler,
	}
	return server.Serve(listener)
}

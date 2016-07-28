package socket

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

// ListenAndServeHTTP starts a HTTP server which is binded to the specified socket path.
func ListenAndServeHTTP(path string, handler http.Handler) error {
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

// GetUniquePath returns path that has nothing on it.
func GetUniquePath(prefix string) (string, error) {
	tmpFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	return tmpFile.Name(), nil
}

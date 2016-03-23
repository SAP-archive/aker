package socket

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
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

func GetUniqueSocketPath(prefix string) (string, error) {
	tmpFile, err := ioutil.TempFile("", prefix)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	return tmpFile.Name(), nil
}

package socket

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sync"
)

type HTTPServer struct {
	path         string
	listener     *net.UnixListener
	requestGroup sync.WaitGroup
	handler      http.Handler
}

func NewHTTPServer(path string, h http.Handler) *HTTPServer {
	return &HTTPServer{
		path:    path,
		handler: h,
	}
}

// Start starts a HTTP server that listens on the configured socket path.
func (s *HTTPServer) Start() error {
	if s.listener != nil {
		return nil
	}

	var err error
	s.listener, err = net.ListenUnix("unix", &net.UnixAddr{
		Name: s.path,
		Net:  "unix",
	})
	if err != nil {
		return err
	}

	go http.Serve(s.listener, http.HandlerFunc(s.serveHTTP))
	return nil
}

func (s *HTTPServer) serveHTTP(w http.ResponseWriter, req *http.Request) {
	s.requestGroup.Add(1)
	defer s.requestGroup.Done()

	s.handler.ServeHTTP(w, req)
}

// Stop releases all resources allocated by the server.
// It also takes care of removing the socket file from the file system.
func (s *HTTPServer) Stop() error {
	if s.listener == nil {
		return nil
	}

	if err := s.listener.Close(); err != nil {
		return err
	}
	s.listener = nil

	s.requestGroup.Wait()
	return nil
}

// // ListenAndServeHTTP starts a HTTP server which is binded to the specified socket path.
// func ListenAndServeHTTP(path string, handler http.Handler) error {
// 	listener, err := net.Listen("unix", path)
// 	if err != nil {
// 		return err
// 	}
// 	defer listener.Close()
//
// 	server := &http.Server{
// 		Handler: handler,
// 	}
// 	return server.Serve(listener)
// }

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

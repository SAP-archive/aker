package socket_test

import (
	"io"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSocket(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Socket Suite")
}

var EchoHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	io.Copy(w, req.Body)
})

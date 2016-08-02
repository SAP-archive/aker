package socket_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	. "github.infra.hana.ondemand.com/I061150/aker/socket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Proxy", func() {
	var path string
	var socketServer *HTTPServer

	var proxyHandler http.Handler

	JustBeforeEach(func() {
		proxyHandler = ProxyHTTP(path)
	})

	BeforeEach(func() {
		var err error
		path, err = GetUniquePath("aker-test")
		Ω(err).ShouldNot(HaveOccurred())

		socketServer = NewHTTPServer(path, EchoHandler)
		Ω(socketServer.Start()).Should(Succeed())
	})

	AfterEach(func() {
		Ω(socketServer.Stop()).Should(Succeed())
	})

	It("should proxy received HTTP requests to the specified socket path", func() {
		payload := []byte("payload")
		req, err := http.NewRequest("POST", "http://whatsoever", bytes.NewBuffer(payload))
		Ω(err).ShouldNot(HaveOccurred())

		respRecorder := httptest.NewRecorder()
		proxyHandler.ServeHTTP(respRecorder, req)
		// validate that EchoHandler has been called
		Ω(respRecorder.Body.Bytes()).Should(Equal(payload))
	})

})

package plugin_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"

	. "github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin/pluginfakes"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListenAndServeHTTP", func() {

	var factory HandlerFactory
	var config []byte

	var fakeSocket *pluginfakes.FakeSocket
	var fakeNotifier *pluginfakes.FakeNotifier

	var server *Server
	var err error

	BeforeEach(func() {
		fakeSocket = new(pluginfakes.FakeSocket)
		fakeNotifier = new(pluginfakes.FakeNotifier)
	})

	JustBeforeEach(func() {
		server = NewServer(bytes.NewBuffer(config), gologger.DefaultLogger, fakeSocket, fakeNotifier)
		err = server.ListenAndServeHTTP(factory)
	})

	Context("when the config passed to stdin is malformed", func() {

		BeforeEach(func() {
			config = []byte("invalid json")
		})

		It("should return a ConfigDecodeError", func() {
			Ω(err).Should(HaveOccurred())
			_, ok := err.(*ConfigDecodeError)
			Ω(ok).Should(BeTrue())
		})
	})

	Context("when the config passed to stdin is valid", func() {

		BeforeEach(func() {
			config = buildConfig("", "")
		})

		Context("and the factory returns an error", func() {
			var factoryErr error

			BeforeEach(func() {
				factoryErr = errors.New("blow!")
				factory = func(_ []byte) (http.Handler, error) {
					return nil, factoryErr
				}
			})

			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
				Ω(err).Should(Equal(factoryErr))
			})
		})

		Context("and the factory returns non-nil handler and no error", func() {

			var socketPath string
			var handler http.Handler
			var httpServer *pluginfakes.FakeHTTPServer

			BeforeEach(func() {
				socketPath = "/tmp/aker.sock"
				config = buildConfig(socketPath, "")

				handler = http.NewServeMux() // we do not care what this is, as long as it is not nil
				factory = func(_ []byte) (http.Handler, error) {
					return handler, nil
				}

				httpServer = new(pluginfakes.FakeHTTPServer)
				fakeSocket.NewHTTPServerReturns(httpServer)
				fakeNotifier.NotifyStub = func(c chan<- os.Signal, sig ...os.Signal) {
					// make sure that ListenAndServeHTTP does not block on this channel
					close(c)
				}
			})

			It("should start HTTP server with the returned handler", func() {
				Ω(fakeSocket.NewHTTPServerCallCount()).Should(Equal(1))
				argSocketPath, argHandler := fakeSocket.NewHTTPServerArgsForCall(0)
				Ω(argSocketPath).Should(Equal(socketPath))
				Ω(argHandler).Should(Equal(handler))

				Ω(httpServer.StartCallCount()).Should(Equal(1))
			})

			It("should stop the HTTP server when signaled", func() {
				Ω(httpServer.StopCallCount()).Should(Equal(1))
			})

			Context("and the config has empty ForwardSocketPath field", func() {
				It("should not call socket.ProxyHTTP", func() {
					Ω(fakeSocket.ProxyHTTPCallCount()).Should(BeZero())
				})
			})

			Context("and the config has non-empty ForwardSocketPath field", func() {

				var forwardSocketPath string

				BeforeEach(func() {
					forwardSocketPath = "non-empty"
					config = buildConfig("", forwardSocketPath)
				})

				It("should create socket HTTP proxy to ForwardSocketPath", func() {
					Ω(fakeSocket.ProxyHTTPCallCount()).Should(Equal(1))
					path := fakeSocket.ProxyHTTPArgsForCall(0)
					Ω(path).Should(Equal(forwardSocketPath))
				})
			})
		})
	})
})

func buildConfig(socketPath, forwardSocketPath string) []byte {
	return []byte(fmt.Sprintf(`{"socket_path":"%s","forward_socket_path":"%s"}`,
		socketPath, forwardSocketPath))
}

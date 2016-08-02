package socket_test

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	. "github.infra.hana.ondemand.com/I061150/aker/socket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	Describe("HTTPServer", func() {

		var path string
		var handler http.Handler
		var server *HTTPServer

		var err error

		JustBeforeEach(func() {
			server = NewHTTPServer(path, handler)
			err = server.Start()
		})

		AfterEach(func() {
			Ω(server.Stop()).Should(Succeed())
		})

		Context("when the server is started on invalid path", func() {
			const invalidPath = "/tmp/00000000001111111111222222222233333333334444444444555555555566666666667777777777888888888899999999.sock"

			BeforeEach(func() {
				path = invalidPath
			})

			It("should return an error", func() {
				Ω(err).Should(HaveOccurred())
			})

		})

		Context("when the server is started on valid path", func() {

			BeforeEach(func() {
				path, err = GetUniquePath("aker-test")
				Ω(err).ShouldNot(HaveOccurred())
				handler = EchoHandler
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should serve requests by the passed http.Handler", func() {
				payload := []byte("payload")

				http := socketHTTPClient(path)
				resp, err := http.Post("http://whatsoever", "text/plain", bytes.NewBuffer(payload))
				Ω(err).ShouldNot(HaveOccurred())
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				Ω(err).ShouldNot(HaveOccurred())

				Ω(body).Should(Equal(payload))
			})

			Context("and is then stopped", func() {
				It("should clean up the socket file", func() {
					server.Stop()
					_, err := os.Stat(path)
					Ω(os.IsNotExist(err)).Should(BeTrue())
				})
			})
		})
	})

	Describe("GetUniquePath", func() {

		var path string
		var err error

		BeforeEach(func() {
			path, err = GetUniquePath("prefix")
		})

		It("should return a path that has nothing on it", func() {
			Ω(err).ShouldNot(HaveOccurred())
			_, statErr := os.Stat(path)
			Ω(statErr).Should(HaveOccurred())
			Ω(os.IsNotExist(statErr)).Should(BeTrue())
		})

	})
})

func socketHTTPClient(path string) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(_, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		},
	}
}

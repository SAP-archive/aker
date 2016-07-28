package socket_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"syscall"

	. "github.infra.hana.ondemand.com/I061150/aker/socket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Socket", func() {
	var socketPath string

	BeforeEach(func() {
		var err error
		socketPath, err = GetUniquePath("aker-test")
		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		err := os.Remove(socketPath)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("when listening on a given socket path", func() {
		var serverCmd *exec.Cmd

		BeforeEach(func() {
			serverCmd = exec.Command("./" + testServerName)
			serverCmd.Stdin = strings.NewReader(socketPath)
			serverCmd.Stdout = os.Stdout
			serverCmd.Stderr = os.Stderr
			err := serverCmd.Start()
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func(done Done) {
			err := serverCmd.Process.Signal(syscall.SIGTERM)
			Ω(err).ShouldNot(HaveOccurred())
			serverCmd.Wait()
			close(done)
		}, 5)

		It("is possible to connect to that server", func() {
			request, err := http.NewRequest("GET", "http://doesnot/matter", nil)
			Ω(err).ShouldNot(HaveOccurred())
			response := httptest.NewRecorder()
			response.Code = -1 // Make sure we have been actually called

			handler := ProxyHTTP(socketPath)
			handler.ServeHTTP(response, request)

			Ω(response.Code).Should(Equal(http.StatusOK))
			Ω(response.Body.String()).Should(Equal("You shall pass, this time!"))
		})
	})
})

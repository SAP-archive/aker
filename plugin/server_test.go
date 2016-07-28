package plugin_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	. "github.infra.hana.ondemand.com/I061150/aker/plugin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListenAndServeHTTP", func() {

	var stdin *os.File
	var factory HandlerFactory
	var err error

	BeforeEach(func() {
		fd, fErr := ioutil.TempFile("", "aker-test")
		Ω(fErr).ShouldNot(HaveOccurred())
		stdin = fd
	})

	AfterEach(func() {
		stdin.Close()
		os.Remove(stdin.Name())
	})

	JustBeforeEach(func() {
		stdin.Seek(0, 0)
		os.Stdin = stdin
		err = ListenAndServeHTTP(factory)
	})

	Context("when the config passed to stdin is malformed", func() {
		BeforeEach(func() {
			stdin.Write([]byte("invalid json"))
		})

		It("should return a ConfigDecodeError", func() {
			Ω(err).Should(HaveOccurred())
			_, ok := err.(*ConfigDecodeError)
			Ω(ok).Should(BeTrue())
		})
	})

	Context("when the config passed to stdin is valid", func() {
		BeforeEach(func() {
			stdin.Write([]byte("{}\n"))
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
	})

})

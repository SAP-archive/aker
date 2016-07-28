package socket_test

import (
	"os"

	. "github.infra.hana.ondemand.com/I061150/aker/socket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
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

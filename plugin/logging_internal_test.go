package plugin

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogWriter", func() {
	var data []byte
	var logger *logWriter
	var sink *bytes.Buffer

	var n int
	var err error

	BeforeEach(func() {
		data = []byte("two\nlines\n")
		sink = new(bytes.Buffer)
		logger = newLogWriter("name", sink)
		n, err = logger.Write(data)
	})

	It("should have returned the number of bytes written", func() {
		Ω(n).Should(Equal(len(data)))
	})

	It("should have not returned an error", func() {
		Ω(err).ShouldNot(HaveOccurred())
	})

	It("should have properly formatted the lines to the sink", func() {
		Ω(sink.String()).Should(Equal("[name]: two\n[name]: lines\n"))
	})

})

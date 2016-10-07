package logging_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/SAP/aker/logging"
	"github.com/SAP/aker/logging/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type countHandler struct {
	serveCount int
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("Sample"))
	Ω(err).ShouldNot(HaveOccurred())
	h.serveCount++
}

type countTimeProvider struct {
	nowCount int
}

func (t *countTimeProvider) Now() time.Time {
	t.nowCount++
	return time.Time{}
}

func (t *countTimeProvider) NowCallCount() int { return t.nowCount }

var _ = Describe("LoggingHandler", func() {

	var timeProvider *countTimeProvider
	var formatter *fakes.FakeFormatter

	var out *bytes.Buffer
	var handler http.Handler

	BeforeEach(func() {
		timeProvider = new(countTimeProvider)
		formatter = new(fakes.FakeFormatter)
		formatter.FormatReturns("formated")
		factory := HandlerFactory{
			TimeProvider: timeProvider.Now,
			Formatter:    formatter,
		}

		out = new(bytes.Buffer)
		handler = factory.LoggingHandler(out, &countHandler{0})
	})

	Context("when request finishes", func() {
		var req *http.Request
		var resp *httptest.ResponseRecorder

		BeforeEach(func() {
			req = request("GET", "http://hack.me", "")
			resp = new(httptest.ResponseRecorder)
			handler.ServeHTTP(resp, req)
		})

		It("should have sampled start and finish timestamps", func() {
			Ω(timeProvider.NowCallCount()).Should(Equal(2))
		})

		It("should have called the formatter", func() {
			Ω(formatter.FormatCallCount()).Should(Equal(1))

			accessEntryArg := formatter.FormatArgsForCall(0)

			Ω(accessEntryArg.StartedAt).Should(Equal(time.Time{}))
			Ω(accessEntryArg.FinishedAt).Should(Equal(time.Time{}))

			Ω(accessEntryArg.Request).Should(Equal(req))
		})

		It("should have written to output sink", func() {
			Ω(out.String()).Should(Equal("formated"))
		})
	})

})

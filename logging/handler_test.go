package logging_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.infra.hana.ondemand.com/I061150/aker/logging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type dummyHandler struct{}

func (h dummyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("Sample"))
	Ω(err).ShouldNot(HaveOccurred())
}

var _ = Describe("LoggingHandler", func() {

	const url = "http://someurl.com/path"

	var log *bytes.Buffer
	var handler http.Handler

	BeforeEach(func() {
		log = new(bytes.Buffer)
		handlerFactory := LoggingHandlerFactory{
			TimeProvider: func() time.Time {
				return time.Time{}
			},
		}
		handler = handlerFactory.LoggingHandler(log, dummyHandler{})
	})

	DescribeTable("access entry logging", func(req *http.Request, expected string) {
		log.Reset()
		w := &httptest.ResponseRecorder{
			Body: new(bytes.Buffer),
		}
		handler.ServeHTTP(w, req)
		Ω(log.String()).Should(Equal(expected))
	},
		Entry("GET", newRequest("GET", url, ""),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "GET /path HTTP/1.1" 200 0 6 "-" "-"  aker_request_id:- response_time:0.000000`)),
		Entry("POST", newRequest("POST", url, "data"),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "POST /path HTTP/1.1" 200 4 6 "-" "-"  aker_request_id:- response_time:0.000000`)),
		Entry("DELETE", newRequest("DELETE", url, "delete data"),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "DELETE /path HTTP/1.1" 200 11 6 "-" "-"  aker_request_id:- response_time:0.000000`)),
		Entry("GET with query params", newRequest("GET", url+"?q=1", ""),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "GET /path?q=1 HTTP/1.1" 200 0 6 "-" "-"  aker_request_id:- response_time:0.000000`)),
		Entry("X-Aker-Request-Id", newRequestWithHeader("GET", url, "", "X-Aker-Request-Id:15"),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "GET /path HTTP/1.1" 200 0 6 "-" "-"  aker_request_id:15 response_time:0.000000`)),
		Entry("User-Agent", newRequestWithHeader("GET", url, "", "User-Agent:godzilla"),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "GET /path HTTP/1.1" 200 0 6 "-" "godzilla"  aker_request_id:- response_time:0.000000`)),
		Entry("Referer", newRequestWithHeader("GET", url, "", "Referer:ref"),
			line(`someurl.com - [01/01/0001:00:00:00 +0000] "GET /path HTTP/1.1" 200 0 6 "ref" "-"  aker_request_id:- response_time:0.000000`)),
	)
})

func line(s string) string {
	return s + "\n"
}

func newRequest(method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(body))
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func newRequestWithHeader(method, url, body, header string) *http.Request {
	req := newRequest(method, url, body)
	kv := strings.Split(header, ":")
	if len(kv) != 2 {
		log.Fatalf("Unsupported header format %q\n", header)
	}
	req.Header.Set(kv[0], kv[1])
	return req
}

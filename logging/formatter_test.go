package logging_test

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/SAP/aker/logging"
	"github.com/SAP/aker/logging/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultFormatter", func() {

	const url = "http://test.me/v2"

	var formatter Formatter
	var now time.Time

	BeforeEach(func() {
		formatter = DefaultFormatter
		now = time.Now()
	})

	DescribeTable("DefaultFormatter", func(entry *AccessEntry, expected string) {
		Ω(formatter.Format(entry)).Should(Equal(expected))
	},
		Entry("GET",
			&AccessEntry{
				Request:    request("GET", url, ""),
				Response:   response(http.StatusOK, 0),
				StartedAt:  now,
				FinishedAt: now.Add(time.Second),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "GET /v2 HTTP/1.1" 200 0 0 "-" "-"  aker_request_id:- response_time:1.000000`+"\n"),
		Entry("GET with query params",
			&AccessEntry{
				Request:    request("GET", url+"?q=1", ""),
				Response:   response(http.StatusOK, 0),
				StartedAt:  now,
				FinishedAt: now.Add(time.Second),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "GET /v2?q=1 HTTP/1.1" 200 0 0 "-" "-"  aker_request_id:- response_time:1.000000`+"\n"),
		Entry("POST",
			&AccessEntry{
				Request:    request("POST", url, "{}"),
				Response:   response(http.StatusCreated, 42),
				StartedAt:  now,
				FinishedAt: now.Add(time.Millisecond),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "POST /v2 HTTP/1.1" 201 2 42 "-" "-"  aker_request_id:- response_time:0.001000`+"\n"),
		Entry("X-Aker-Request-Id",
			&AccessEntry{
				Request:    requestWithHeader("GET", url, "", "X-Aker-Request-Id:15"),
				Response:   response(http.StatusOK, 42),
				StartedAt:  now,
				FinishedAt: now.Add(time.Second),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "GET /v2 HTTP/1.1" 200 0 42 "-" "-"  aker_request_id:15 response_time:1.000000`+"\n"),
		Entry("User-Agent",
			&AccessEntry{
				Request:    requestWithHeader("GET", url, "", "User-Agent:mozilla"),
				Response:   response(http.StatusOK, 42),
				StartedAt:  now,
				FinishedAt: now.Add(time.Second),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "GET /v2 HTTP/1.1" 200 0 42 "-" "mozilla"  aker_request_id:- response_time:1.000000`+"\n"),
		Entry("Referer",
			&AccessEntry{
				Request:    requestWithHeader("GET", url, "", "Referer:me"),
				Response:   response(http.StatusOK, 42),
				StartedAt:  now,
				FinishedAt: now.Add(time.Second),
			},
			`test.me - [01/01/0001:00:00:00 +0000] "GET /v2 HTTP/1.1" 200 0 42 "me" "-"  aker_request_id:- response_time:1.000000`+"\n"),
	)
})

func request(method, path, body string) *http.Request {
	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	if err != nil {
		log.Fatalf("error creating request: %v", err)
	}
	return req
}

func requestWithHeader(method, path, body, header string) *http.Request {
	req := request(method, path, body)
	kv := strings.Split(header, ":")
	if len(kv) != 2 {
		log.Fatalf("unsupported header format: %q\n", header)
	}
	req.Header.Add(kv[0], kv[1])
	return req
}

func response(status, size int) ResponseRecorder {
	resp := new(fakes.FakeResponseRecorder)
	resp.StatusReturns(status)
	resp.SizeReturns(size)
	return resp
}

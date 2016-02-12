package adapter

import (
	"mime/multipart"
	"net/http"
	"net/url"

	"github.wdf.sap.corp/I061150/aker/api"
)

func NewRequest(delegate api.Request) *http.Request {
	return &http.Request{
		URL:           delegate.URL(),
		Method:        delegate.Method(),
		Host:          delegate.Host(),
		Body:          delegate,
		ContentLength: int64(delegate.ContentLength()),
		Header:        delegate.Headers(),
		Proto:         "HTTP/1.1", // FIXME
		ProtoMajor:    1,          // FIXME
		ProtoMinor:    1,          // FIXME
		RemoteAddr:    "",         // FIXME
		RequestURI:    "/",        // FIXME
		Close:         false,      // FIXME

		TransferEncoding: []string{},
		Form:             url.Values{},
		PostForm:         url.Values{},
		MultipartForm:    &multipart.Form{},
		Trailer:          make(map[string][]string),
	}
}

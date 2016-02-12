package api

import "net/url"

type Request interface {
	URL() *url.URL
	Method() string
	Host() string
	ContentLength() int
	Headers() map[string][]string
	Header(name string) string
	Read([]byte) (int, error)
	Close() error
}

type Response interface {
	SetHeader(name string, values []string)
	WriteStatus(int)
	Write([]byte) (int, error)
}

type Data interface {
	SetString(name, value string)
	String(name string) string
}

type Context struct {
	Request  Request
	Response Response
	Data     Data
}

type Plugin interface {
	Config([]byte) error
	Process(Context) bool
}

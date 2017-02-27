/*
Pakcage plugin provides means to create an Aker plugin. It exposes the
InitFunc type, which handles the initialization of plugins.

Following is an example plugin implementation that returns configured
response body and status code.

	package main

	import (
		"net/http"

		"github.com/SAP/aker/plugin"
	)

	type MessageHandler struct {
		Message []byte
		Code    int
	}

	func (h *MessageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(h.Code)
		w.Write(h.Message)
	}

	func Init(config []byte) (http.Handler, error) {
		var handler = &MessageHandler{}
		if err := plugin.UnmarshalConfig(data, &handler); err != nil {
			return nil, err
		}
		return handler, nil
	}

Creating a plugin requires implementing and exporting an Init function, which
basically is just a function that accepts a byte array as input parameter
and returns a http.Handler and an error. The byte array provided as input
contains all configuration that was passed to the Aker process and is
intended to configure the implemented plugin.

The Init func should use UnmarshalConfig to get a Go struct representation
of the data. The underlying notation is YAML, thus if you want to make use
of Go's struct field tags, you should use 'yaml:<opts>' for configuring
how the unmarshaller should handle the Go struct fields.

	func Init(data []byte) (http.Handler, error) {
		var myConfig struct{
			a int
			b string
		} cfg

		if err := plugin.UnmarshalConfig(data, &cfg); err != nil {
			return nil, err
		}
		return newMyHandler(cfg), nil
	}

The plugin is loaded in memory through Go's built-in plugin functionality.
The plugin should be built using:
go build -buildmode=plugin
So it is possible to load it as dynamic library. The main package of the
plugin should export a Init sybmol, which has the plugin.InitFunc signature.

Since loading the plugin makes it part of the aker process, writing to
stdout and stderr is the way to send log messages to the central Aker log.

Request tracking mechanism is provided by Aker. Each incoming HTTP request
is decorated with the X-Aker-Request-Id header, as well as each response.
If a plugin encounters problem with some request, it is advisable to dump
the value of the X-Aker-Request-Id header for debugging purposes.
The X-Aker-Request-Id header is also propagated to the end user, so when
one complains, she is able to provide the header value for tracing.

If a plugin is part of a plugin chain, which means that each request gets
processed by multiple plugins before it is returned to Aker and thus to the
user, then the way of telling the requests not to continue further the plugin
chain is to write something to the response by calling Write or WriteHeader of
the http.Responsewriter. This will stop the request from going through
sequential plugins and will return the response to the end user.
*/
package plugin

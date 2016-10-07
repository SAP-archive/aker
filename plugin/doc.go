/*
  Pakcage plugin provides means to create an Aker plugin. It exposes the
  ListenAndServeHTTP function, which is the main entry point for a plugin.

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

    func main() {
    	plugin.ListenAndServeHTTP(func(data []byte) (http.Handler, error) {
    		var handler = &MessageHandler{}
    		if err := plugin.UnmarshalConfig(data, &handler); err != nil {
    			return nil, err
    		}
    		return handler, nil
    	})
    }

  Creating a plugin requires implementing a HandlerFactory object, which
  basically is just a function that accepts a byte array as input parameter
  and returns a http.Handler and an error. The byte array provided as input
  contains all configuration that was passed to the Aker process and is
  intended to configure the implemented plugin.

  The factory should use UnmarshalConfig to get a Go struct representation
  of the data. The underlying notation is YAML, thus if you want to make use
  of Go's struct field tags, you should use 'yaml:<opts>' for configuring
  how the unmarshaller should handle the Go struct fields.

  	func myFactory(data []byte) (http.Handler, error) {
  		var myConfig struct{
  			a int
  			b string
  		} cfg

  		if err := plugin.UnmarshalConfig(data, &cfg); err != nil {
  			return nil, err
  		}
  		return newMyHandler(cfg), nil
  	}

  The communication between Aker and each plugin, and between each pair of
  plugins, happens via HTTP, which is transported over unix domain sockets.
  The ListenAndServeHTTP method takes care of cleaning up the socket file,
  once the plugin receives signal to exit. Because of that, it is undesirable
  to call os.Exit from within a plugin, since this will leave the allocated
  socket file on the file system.

  Plugin's Stdin and Stderr are captured by Aker, so writing to them is the
  way to send log messages to the central Aker log. They'll get decorated by
  appending the plugin name in front of each log line.

  	[plugin-name]: 2016/07/19 14:28:59 [INFO] Starting...

  It is advisable to use the logger provided by the logging package.

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

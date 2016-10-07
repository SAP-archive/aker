Aker
====

Aker is an HTTP reverse proxy server that is designed to handle incoming requests by forwarding them to a number of plugin handlers.

Plugins can be configured and attached to specific HTTP paths. By default, an HTTP `Not Found` status is returned for paths that don't have plugins configured.

## License
This project is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.

## Getting Started

Before setting up Aker, one should get familiar with the internal workings of the software.

Aker acts as an HTTP reverse proxy that routes incoming requests to sequences of plugins, depending on the request path.
All of the plugin chains (sequences) and their path binding are modeled via the Aker configuration.

![Overview](https://github.infra.hana.ondemand.com/raw/cloudfoundry/aker/master/documents/overview.png)

Once a request is received by Aker, it checks the path of the request and determines the plugin chain that should handle it. Afterwards, it forwards the request to the first plugin in that chain.
That plugin, external to the Aker process, has the chance to process and modify the incoming request. It can then write output to the response, at which point the processing stops and the request will not go further, or it can leave the request to be processed by subsequent plugins in the chain.

For example, there could be a plugin which does basic authentication. It checks the `Authorization` header in the request and if the user and password match some desired value, then control is passed to the next plugin, otherwise a `401 Unauthorized` code can be returned and processing flow stopped.

In the example above, we have two plugin chains. The first one processes requests to the `/hello.png` endpoint by passing the request to `Plugin A` and potentially to `Plugin B` if the former allows it. The second plugin chain, consisting of a single plugin, processes requests to the `/auth` path by passing control to `Plugin C`.

An endpoint can handle a path tree as well. If in the diagram above the second plugin had the `/auth/` path configured instead, the handler would handle `/auth/something/else` as well as the root `/auth` path.
If two endpoints have overlapping paths, then the plugin chain with a path expression that matches the incoming request better will be the one to handle it. This is aligned with how the Go language processes requests.

A common use case for Aker is to have a `/` endpoint that has a plugin chain ending with a reverse-proxy plugin. That way, anything that does not get processed by other more specific `chains` gets reverse proxied.

Plugins can use HTTP headers to forward information to subsequent plugins in the chain.

## User Guide

The application is written in Go so you will need to set that up. Once you have Go, you can use the following command to download the source code and build it.

```bash
go get github.com/SAP/aker
```

To verify that Aker has been properly installed, use the following command.

```bash
which aker
```

You can run Aker as follows.

```bash
aker -config <path_to_config_file>
```

:information_source: If you don't specify the `-config` flag, then Aker will look for a configuration file in `./config.yml`.

Let's have a look at a minimal configuration.

```yaml
---
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 5
  write_timeout: 10

endpoints: []
```

This configuration will start Aker and have it listen on the `8080` port for HTTP requests and will have the corresponding read and write timeouts for tcp connections. It will not handle any of the requests, however, since we haven't configured any behavior.

:information_source: If you want Aker to listen only for local requests, you can change `host` from `0.0.0.0` to `127.0.0.1`.

Here is an extension to the above configuration that adds some meaningful behavior.

```yaml
---
server:
  host: 0.0.0.0
  port: 8080
  read_timeout: 5
  write_timeout: 10

endpoints:
  - path: "/"
    audit: true
    plugins:
      - name: aker-proxy-plugin
        configuration:
          url: http://location.com
```

The above configuration specifies that all incoming requests on path `/` should be forwarded to the [aker-proxy-plugin](https://github.com/SAP/aker-proxy-plugin), which in turn will proxy calls to `http://example.org`.

:information_source: One needs to make sure that the `aker-proxy-plugin` plugin is available on the `PATH`, or one could configure the plugin `name` to point to the plugin executable.

The `audit` option can be used to configure detailed logging of incoming requests.

## Developer Guide

You will need to download the following tools.

* [Ginkgo](https://github.com/onsi/ginkgo) - Used for running the tests.
* [Godep](https://github.com/tools/godep) - Used for dependency management.
* [Counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) - Used for interface mocking.

You can run the tests with the following command.

```bash
ginkgo -r
```

You can regenerate all the interface mocks with the following command.

```bash
go generate ./...
```

You can update the stored Go dependencies with the following command.

```
rm -rf Godeps/ vendor/
godep save ./...
```

## Writing a Plugin

Package `plugin` of Aker provides means to create an Aker plugin. It exposes the `ListenAndServeHTTP` function, which is the main entry point for a plugin.

Following is an example plugin implementation that returns a response body and status code, both specified via the Aker configuration.

```go
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
```

Creating a plugin requires implementing a `HandlerFactory` object, which is basically just a function that accepts a byte array as input parameter and returns a http.Handler and an error.
The byte array will contain the plugin configuration as specified in the Aker `yaml` configuration.

The factory should use `UnmarshalConfig` to get a Go struct representation of the data. The underlying notation is YAML, thus if you want to make use of Go's struct field tags, you should use `yaml:<opts>` for configuring
how the unmarshaller should handle the Go struct fields.

```go
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
```

The communication between Aker and each plugin, and between each pair of plugins, happens via HTTP, which is transported over unix domain sockets.
The `ListenAndServeHTTP` function takes care of cleaning up the socket file, once the plugin receives a signal to exit. Because of that, it is undesirable to call `os.Exit` from within a plugin, as this will leave the allocated socket file on the file system.

Plugin's `stdout` and `stderr` are captured by Aker, so writing to them is the way to send log messages to the central Aker log. They'll get decorated by having the plugin name appended in front of each log line.

```
[plugin-name]: 2016/07/19 14:28:59 [INFO] Starting...
```

It is advisable to use the logger provided by the `logging` package.

Request tracking mechanism is provided by Aker. Each incoming HTTP request is decorated with the `X-Aker-Request-Id` header, as well as each response.
If a plugin encounters problem with some request, it is advisable to dump the value of the `X-Aker-Request-Id` header for debugging purposes.
The `X-Aker-Request-Id` header is also propagated to the end user, so should the user face a problem, they can provide the header value for tracing.

If a plugin is part of a plugin chain, which means that each request gets processed by multiple plugins before it is returned to Aker and thus to the user, then the way of telling the requests not to continue further the plugin chain is to write something to the response by calling Write or WriteHeader of the `http.ResponseWriter`. This will stop the request from going through subsequent plugins and will return the response to the end user.

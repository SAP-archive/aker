// Used only for testing.
package main

import (
	"log"
	"net/http"

	"github.infra.hana.ondemand.com/I061150/aker/plugin"
)

type configHandler []byte

func (h configHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Write(h)
}

func configHandlerFactory(config []byte) (http.Handler, error) {
	return configHandler(config), nil
}

func main() {
	err := plugin.ListenAndServeHTTP(configHandlerFactory)
	if err != nil {
		log.Fatal(err)
	}
}

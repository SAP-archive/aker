package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.infra.hana.ondemand.com/I061150/aker/socket"
)

const greeting = "You shall pass, this time!"

func main() {
	socketPath := readSocketPath()
	socket.ListenAndServeHTTP(socketPath, http.HandlerFunc(serverFunc))
}

func readSocketPath() string {
	socketPath, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(socketPath)
}

func serverFunc(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, greeting)
}

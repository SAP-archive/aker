package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.wdf.sap.corp/I061150/aker/socket"
)

const greeting = "You shall pass, this time!"

func main() {
	socketPath := readSocketPath()
	socket.ListenAndServe(socketPath, http.HandlerFunc(serverFunc))
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

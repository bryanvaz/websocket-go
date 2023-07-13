package main

import (
	"flag"
	"fmt"
	// "log"
	"bryanvaz/wss/pkg/server"
	"net/http"
)

// var addr = flag.String("addr", "0.0.0.0:42069", "http service address")

func main() {

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Setting up the server!")
	})

	wsServer, err := server.NewServer()
	if err != nil {
		fmt.Println("Error creating server: ", err)
	}
	http.HandleFunc(("/ws"), wsServer.HandleNewWsConn)
	http.ListenAndServe(":8080", nil)

}

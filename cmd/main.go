package main

import (
	"flag"
	"fmt"
	"log"

	// "log"
	"bryanvaz/wss/pkg/server"
	"net/http"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func main() {

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "visit /ws for websocket connection")
	})

	wsServer, err := server.NewServer()
	if err != nil {
		fmt.Println("Error creating server: ", err)
	}
	http.HandleFunc(("/ws"), wsServer.HandleNewWsConn)
	log.Println("Starting up server on " + *addr + "...")
	http.ListenAndServe(*addr, nil)
}

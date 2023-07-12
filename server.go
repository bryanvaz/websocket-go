package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// var addr = flag.String("addr", "0.0.0.0:42069", "http service address")

// type Server struct {
// }

// func NewServer() {

// }

var upgrader = websocket.Upgrader{}

func main() {

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Setting up the server!")
	})
	http.HandleFunc(("/ws"), func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(`WS Upgrade Failed: ${err}`)
		}
		defer conn.Close()

		// Read message from client
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(`Msg read failed: ${err}`)
				return
			}
			log.Printf("Message Received: %s\n", msg)

			// Echo message back to client wich cat
			echoResponse := fmt.Sprintf("*meow* %s", msg)
			err = conn.WriteMessage(websocket.TextMessage, []byte(echoResponse))
			if err != nil {
				log.Println(`Write Failed: ${err}`)
				return
			} else {
				log.Printf("Message Sent: %s\n", echoResponse)
			}

		}
	})
	http.ListenAndServe(":8080", nil)

}

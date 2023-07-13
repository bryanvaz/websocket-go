package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct{}

var upgrader = websocket.Upgrader{}

func NewServer() (*Server, error) {
	server := Server{}

	return &server, nil
}

func respGenerator(input string) string {
	return fmt.Sprintf("*meow* %s", input)
}

func (s *Server) HandleNewWsConn(w http.ResponseWriter, r *http.Request) {
	log.Println("New WS Connection requested from: " + r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(`WS Upgrade Failed: ${err}`)
	}
	defer conn.Close()

	log.Println("WS Connection established with: " + r.RemoteAddr)

	// Read message from client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("RECV FAIL: " + err.Error())
			return
		}
		log.Printf("RECV (%s): %s\n", r.RemoteAddr, msg)

		// Echo message back to client wich cat
		echoResponse := respGenerator((string(msg)))
		err = conn.WriteMessage(websocket.TextMessage, []byte(echoResponse))
		if err != nil {
			log.Println("SEND FAIL: " + err.Error())
			return
		} else {
			log.Printf("SEND (%s): %s\n", r.RemoteAddr, echoResponse)
		}

	}
}

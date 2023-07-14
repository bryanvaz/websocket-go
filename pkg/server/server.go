package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	clientList map[*ClientConnection]bool
	mutex      sync.Mutex
}

func NewServer() (*Server, error) {
	server := Server{
		clientList: make(map[*ClientConnection]bool),
	}

	return &server, nil
}

func respGenerator(input string) string {
	return fmt.Sprintf("*meow* %s", input)
}

func (s *Server) AddClient(client *ClientConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clientList[client] = true
}

func (s *Server) RemoveClient(client *ClientConnection) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.clientList, client)
	log.Println("Client disconnected: ", client.conn.RemoteAddr())
}

func (s *Server) HandleNewWsConn(w http.ResponseWriter, r *http.Request) {

	client, err := NewClientConnection(w, r, s)
	if err != nil {
		log.Println("Error creating new client connection: ", err)
	}

	s.AddClient(client)
	log.Println("New client connected: ", r.RemoteAddr)
	log.Println("Total clients connected: ", len(s.clientList))

	go func() {
		for {
			msg := <-client.inBound
			// client.outBound <- respGenerator(msg)
			for destClient := range s.clientList {
				destClient.outBound <- respGenerator(msg)
			}
		}
	}()

}

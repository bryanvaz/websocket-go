package server

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct{}

func NewServer() (*Server, error) {
	server := Server{}

	return &server, nil
}

func respGenerator(input string) string {
	return fmt.Sprintf("*meow* %s", input)
}

func (s *Server) HandleNewWsConn(w http.ResponseWriter, r *http.Request) {

	_, err := NewClientConnection(w, r)
	if err != nil {
		log.Println("Error creating new client connection: ", err)
	}

}

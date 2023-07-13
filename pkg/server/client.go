package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

// Represents the connection of a single client
type ClientConnection struct{}

// Upgrades to WS connection for an HTTP request
func NewClientConnection(w http.ResponseWriter, r *http.Request) (*ClientConnection, error) {
	log.Printf("WS UPGRADE REQUEST (%s)\n", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("WS UPGRADE FAIL ("+r.RemoteAddr+"): ", err)
		return nil, err
	}
	log.Printf("WS UPGRADE SUCCESS (%s)\n", r.RemoteAddr)

	clientConn := ClientConnection{}

	go func() {
		defer func() {
			conn.Close()
		}()

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("RECV FAILURE (%s): %s\n", r.RemoteAddr, err)
				break
				// If you don't break, the loop will
				// continue to try to read from the connection
			}
			if mt != websocket.TextMessage {
				continue // Ignore non-text messages
			}
			log.Printf("RECV (%s): %s\n", r.RemoteAddr, msg)

			// Echo message back to client as a cat
			resp := respGenerator((string(msg)))
			err = conn.WriteMessage(websocket.TextMessage, []byte(resp))
			if err != nil {
				log.Printf("SEND FAILURE (%s): %s\n", r.RemoteAddr, err)
				break
				// Should probably break because there's something wrong
				// with the connection
			} else {
				log.Printf("SEND (%s): %s\n", r.RemoteAddr, resp)
			}
		}
	}()

	return &clientConn, nil
}

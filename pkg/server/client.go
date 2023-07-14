package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Represents the connection of a single client
type ClientConnection struct {
	conn   *websocket.Conn
	closed bool
	// from me to network (msg received from this client)
	inBound <-chan string
	// from network to me (msg to send to this client)
	outBound chan<- string
}

// Upgrades to WS connection for an HTTP request
func NewClientConnection(
	w http.ResponseWriter,
	r *http.Request,
	s *Server,
) (*ClientConnection, error) {
	log.Printf("WS UPGRADE REQUEST (%s)\n", r.RemoteAddr)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("WS UPGRADE FAIL ("+r.RemoteAddr+"): ", err)
		return nil, err
	}
	log.Printf("WS UPGRADE SUCCESS (%s)\n", r.RemoteAddr)

	// from network to me
	out := make(chan string, 1)

	// from me to network
	in := make(chan string, 1)

	client := ClientConnection{
		conn:     c,
		closed:   false,
		inBound:  in,
		outBound: out,
	}

	go func() {
		defer func() {
			c.Close()
			s.RemoveClient(&client)
			client.closed = true
		}()

		for {
			mt, msg, err := c.ReadMessage()
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

			in <- string(msg)

			// Echo message back to client as a cat
			// resp := respGenerator((string(msg)))
			// err = c.WriteMessage(websocket.TextMessage, []byte(resp))
			// if err != nil {
			// 	log.Printf("SEND FAILURE (%s) (msg: %s): %s\n", r.RemoteAddr, resp, err)
			// 	break
			// } else {
			// 	log.Printf("SEND (%s): %s\n", r.RemoteAddr, resp)
			// }
		}
	}()

	go func() {
		for msg := range out {
			if !client.closed {
				outMsg := "NETWORK: " + msg
				err = c.WriteMessage(websocket.TextMessage, []byte(outMsg))
				if err != nil {
					log.Printf("SEND FAILURE (%s) (msg: %s): %s\n", r.RemoteAddr, outMsg, err)
					break
				} else {
					log.Printf("SEND (%s): %s\n", r.RemoteAddr, outMsg)
				}
			}
		}
	}()

	return &client, nil
}

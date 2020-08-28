package pump

import (
	"encoding/json"
	"log"
	"time"
	//test

	"github.com/gorilla/websocket"

	"T/MessagingService/hub"
	model "T/MessagingService/models"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2024
)

// readPump pumps messages from the websocket connection to the hub.
func ReadPump(s hub.Subscription) {
	c := s.Conn
	h := hub.H
	defer func() {
		h.UnRegister <- s
		c.Ws.Close()
	}()
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg1 model.MessageData
		err = json.Unmarshal(msg, &msg1)

		m := hub.Message{msg1, s.Room}
		h.Broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func WritePump(s *hub.Subscription) {
	c := s.Conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
				c.Ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Ws.WriteJSON(message); err != nil {
				return
			}
		case <-ticker.C:
			c.Ws.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.Ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

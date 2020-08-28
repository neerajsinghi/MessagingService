package client

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	buffer "T/MessagingService/buffer"
	"T/MessagingService/hub"
	model "T/MessagingService/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2024,
	WriteBufferSize: 2024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var h = hub.H

// connection is an middleman between the websocket connection and the hub.

// GetPeople is an httpHandler for route POST /updatepost
// This is the api used for updating data in post
func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomid")
	log.Println(roomID)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	connec := &hub.Connection{Send: make(chan model.MessageData), Ws: ws}

	subs := hub.Subscription{connec, roomID}
	h.Register <- subs
	go buffer.ReadPump(subs)
	go buffer.WritePump(&subs)
}

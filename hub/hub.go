package hub

import (
	"github.com/gorilla/websocket"

	model "T/MessagingService/models"
)

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan model.MessageData
}
type Message struct {
	Data model.MessageData
	Room string
}

type Subscription struct {
	Conn *Connection
	Room string
}

// hub maintains the set of active connections and Broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Rooms map[string]map[*Connection]bool

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	Register chan Subscription

	// UnRegister requests from connections.
	UnRegister chan Subscription
}

var H = Hub{
	Broadcast:  make(chan Message),
	Register:   make(chan Subscription),
	UnRegister: make(chan Subscription),
	Rooms:      make(map[string]map[*Connection]bool),
}

func init() {
	go H.run()
}
func (h *Hub) run() {
	for {
		select {
		case subscribe := <-h.Register:
			connections := h.Rooms[subscribe.Room]
			if connections == nil {
				connections = make(map[*Connection]bool)
				h.Rooms[subscribe.Room] = connections
			}
			h.Rooms[subscribe.Room][subscribe.Conn] = true
		case unsubscribe := <-h.UnRegister:
			connections := h.Rooms[unsubscribe.Room]
			if connections != nil {
				if _, ok := connections[unsubscribe.Conn]; ok {
					delete(connections, unsubscribe.Conn)
					close(unsubscribe.Conn.Send)
					if len(connections) == 0 {
						delete(h.Rooms, unsubscribe.Room)
					}
				}
			}
		case mesg := <-h.Broadcast:
			connections := h.Rooms[mesg.Room]
			for c := range connections {
				select {
				case c.Send <- mesg.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, mesg.Room)
					}
				}
			}
		}
	}
}

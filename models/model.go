package model

import "github.com/gorilla/websocket"

//FeedStruct Data to be added to the database
type MessageData struct {
	Title        string `json:"title,omitempty"`
	Body         string `json:"body,omitempty"`
	TopicID      string `json:"topic_id,omitempty"`
	Userid       string `json:"userid,omitempty"`
	CurrentValue int    `json:"current_value,omitempty"`
}

type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan MessageData
}
type Message struct {
	Data MessageData
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

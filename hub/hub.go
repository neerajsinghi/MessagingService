package hub

import model "T/MessagingService/models"

var H = model.Hub{
	Broadcast:  make(chan model.Message),
	Register:   make(chan model.Subscription),
	UnRegister: make(chan model.Subscription),
	Rooms:      make(map[string]map[*model.Connection]bool),
}

func init() {
	go run(&H)
}
func run(h *model.Hub) {
	for {
		select {
		case subscribe := <-h.Register:
			connections := h.Rooms[subscribe.Room]
			if connections == nil {
				connections = make(map[*model.Connection]bool)
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

package websocket

import "github.com/gorilla/websocket"

type Hub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan Message
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.Clients[conn] = true

		case conn := <-h.Unregister:
			delete(h.Clients, conn)
			conn.Close()

		case msg := <-h.Broadcast:
			for conn := range h.Clients {
				conn.WriteJSON(msg)

			}
		}
	}
}

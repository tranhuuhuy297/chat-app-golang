package ws

import (
	log "github.com/sirupsen/logrus"
)

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			_, existedRoom := h.Rooms[client.RoomID]
			if !existedRoom {
				return
			}

			room := h.Rooms[client.RoomID]
			log.Printf("ws/hub/Run| Existed room: %v", room)
			if _, ok := room.Clients[client.ID]; !ok {
				room.Clients[client.ID] = client
				log.Printf("ws/hub/Run| Client joined: %v", client)
			}
		case client := <-h.Unregister:
			_, existedRoom := h.Rooms[client.RoomID]
			if !existedRoom {
				return
			}

			clients := h.Rooms[client.RoomID].Clients
			_, existedClient := clients[client.ID]
			if !existedClient {
				return
			}

			h.Broadcast <- &Message{
				Content:  "user left the chat",
				RoomID:   client.RoomID,
				Username: client.Username,
			}

			delete(h.Rooms[client.RoomID].Clients, client.ID)
			close(client.Message)
		case message := <-h.Broadcast:
			_, existedRoom := h.Rooms[message.RoomID]
			if !existedRoom {
				return
			}

			for _, client := range h.Rooms[message.RoomID].Clients {
				client.Message <- message
			}
		}
	}
}

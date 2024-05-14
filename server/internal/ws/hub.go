package ws

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
			if _, ok := h.Rooms[client.RoomID]; ok {
				room := h.Rooms[client.RoomID]

				if _, ok := room.Clients[client.ID]; !ok {
					room.Clients[client.ID] = client
				}
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
			if _, ok := h.Rooms[message.RoomID]; ok {
				for _, client := range h.Rooms[message.RoomID].Clients {
					client.Message <- message
				}
			}
		}
	}
}

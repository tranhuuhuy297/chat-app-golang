package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

type Controller struct {
	hub *Hub
}

func NewController(h *Hub) *Controller {
	return &Controller{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (controller *Controller) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	controller.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (controller *Controller) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Query("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")
	log.Printf("ws/ws_controller/JoinRoom| roomId: %s, userId: %s, username %s", roomID, clientID, username)

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	message := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}
	log.Printf("ws/ws_controller/JoinRoom| message: %v", message)

	controller.hub.Register <- client
	controller.hub.Broadcast <- message

	go client.writeMessage()
	client.readMessage(controller.hub)
}

type GetRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (controller *Controller) GetRooms(c *gin.Context) {
	rooms := make([]GetRoomRes, 0)

	for _, room := range controller.hub.Rooms {
		rooms = append(rooms, GetRoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type GetClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (controller *Controller) GetClients(c *gin.Context) {
	var clients []GetClientRes
	roomId := c.Query("roomId")

	if _, ok := controller.hub.Rooms[roomId]; !ok {
		clients = make([]GetClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range controller.hub.Rooms[roomId].Clients {
		clients = append(clients, GetClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

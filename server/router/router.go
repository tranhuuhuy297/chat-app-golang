package router

import (
	md "server/internal/middleware"
	"server/internal/user"
	"server/internal/ws"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Init(userController *user.Controller, wsController *ws.Controller) {
	r = gin.Default()

	r.POST("/sign-up", userController.CreateUser)
	r.POST("/login", userController.Login)
	r.GET("/logout", userController.Logout)

	r.POST("/ws/rooms", md.AuthenticateMiddleware, wsController.CreateRoom)
	r.GET("/ws/rooms", wsController.GetRooms)
	r.GET("/ws/room-joining", wsController.JoinRoom)
	r.GET("/ws/clients", wsController.GetClients)
}

func Start(address string) error {
	return r.Run(address)
}

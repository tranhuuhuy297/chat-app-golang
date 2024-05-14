package router

import (
	"server/internal/user"

	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Init(userController *user.Controller) {
	r = gin.Default()

	r.POST("/sign-up", userController.CreateUser)
	r.POST("/login", userController.Login)
	r.GET("/logout", userController.Logout)
}

func Start(address string) error {
	return r.Run(address)
}

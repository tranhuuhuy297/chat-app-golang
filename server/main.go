package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"server/db"
	"server/internal/user"
	"server/internal/ws"
	"server/router"

	"github.com/joho/godotenv"
)

var ENDPOINT string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	ENDPOINT = os.Getenv("ENDPOINT")
}

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Printf("main| Could not initialize database connection: %s", err)
	}

	// User
	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userController := user.NewController(userService)

	// Websocket
	hub := ws.NewHub()
	wsController := ws.NewController(hub)
	go hub.Run()

	// Router
	router.Init(userController, wsController)
	router.Start(ENDPOINT)
}

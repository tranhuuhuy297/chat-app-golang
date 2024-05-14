package main

import (
	"log"
	"os"
	"server/db"
	"server/internal/user"
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
		log.Fatalf("could not initialize database connection: %s", err)
	}

	userRepository := user.NewRepository(dbConn.GetDB())
	userService := user.NewService(userRepository)
	userController := user.NewController(userService)

	router.Init(userController)
	router.Start(ENDPOINT)
}

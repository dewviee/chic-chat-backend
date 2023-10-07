// main.go

package main

import (
	"chicchat/controllers"
	"chicchat/middlewares"
	"chicchat/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Create a map to store rooms
	rooms := make(map[string]*models.Room)

	// Use the WebSocket handler
	app.Get("/room/:id",
		middlewares.WebSocketMiddleware(rooms),
		middlewares.RoomAuth(),
		controllers.WebSocketHandler(rooms))

	log.Fatal(app.Listen("localhost:8080"))
}

// main.go

package main

import (
	"chicchat/config"
	"chicchat/controllers"
	"chicchat/database"
	"chicchat/middlewares"
	"chicchat/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("./config/.env")

	db := database.ConnectDB(config.GetMySqlDSN())
	database.MigratingDatabase(db)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Create a map to store rooms
	rooms := make(map[string]*models.RoomWebSocket)

	// Use the WebSocket handler
	app.Get("/room/:id",
		middlewares.WebSocketMiddleware(rooms),
		middlewares.RoomAuth(),
		controllers.WebSocketHandler(rooms))

	log.Fatal(app.Listen("localhost:8080"))
}

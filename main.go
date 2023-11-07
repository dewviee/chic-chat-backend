// main.go

package main

import (
	"chicchat/config"
	"chicchat/database"
	"chicchat/models"
	"chicchat/routes"
	"fmt"
	"log"
	"os"

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
		AllowOrigins: "http://localhost:3000,http://chic-chat-frontend:3000",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Create a map to store rooms
	rooms := make(map[string]*models.RoomWebSocket)

	// Register routes
	routes.AuthRoute(app, db)
	routes.Room(app, db, rooms)
	routes.UserRoute(app, db)

	serverHost := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	log.Fatal(app.Listen(serverHost))
}

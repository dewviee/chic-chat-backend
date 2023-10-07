package routes

import (
	"chicchat/controllers"
	"chicchat/middlewares"
	"chicchat/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Room(app *fiber.App, db *gorm.DB, rooms map[string]*models.RoomWebSocket) {
	// Use the WebSocket handler
	app.Get("/room/:id",
		middlewares.WebSocketMiddleware(rooms),
		middlewares.RoomAuth(),
		controllers.WebSocketHandler(rooms))
}

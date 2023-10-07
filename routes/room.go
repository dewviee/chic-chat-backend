package routes

import (
	"chicchat/models"

	"github.com/gofiber/fiber/v2"
)

func Room(app *fiber.App, rooms map[string]*models.Room) {
	// app.Get("/room/:room", controllers.Room(rooms))
	// app.Get("/room/:room/:id", controllers.ConnectToRoom(rooms))
}

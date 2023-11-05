package routes

import (
	"chicchat/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoute(app *fiber.App, db *gorm.DB) {
	app.Put("/user", controllers.UpdateUser(db))
	app.Get("/user/:id", controllers.GetUser(db))
	app.Get("/user/picture/:filename", controllers.GetProfilePicture(db))
	app.Post("/user/picture", controllers.UploadProfilePicture(db))
}

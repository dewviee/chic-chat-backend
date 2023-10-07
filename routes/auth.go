package routes

import (
	"chicchat/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(app *fiber.App, db *gorm.DB) {
	app.Post("/auth/register/email", controllers.CreateUserByEmail(db))
	app.Post("/auth/login/email", controllers.HandlerEmailLogin(db))
}

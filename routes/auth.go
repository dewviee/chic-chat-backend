package routes

import (
	"chicchat/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoute(app *fiber.App, db *gorm.DB) {
	app.Post("/auth/register/email", controllers.RegisterUserByEmail(db))
	app.Post("/auth/login/email", controllers.HandlerUserLogin(db))
}

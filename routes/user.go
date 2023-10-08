package routes

import (
	"chicchat/controllers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserRoute(app *fiber.App, db *gorm.DB) {
	app.Get("/user/picture/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.SendFile("./assets/pictures/" + id + ".png")
	})
	app.Put("/user", controllers.UpdateUser(db))
}

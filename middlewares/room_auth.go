package middlewares

import "github.com/gofiber/fiber/v2"

func RoomAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID := c.Params("id")

		//TODO: Replace with db query
		if roomID == "1" {
			return c.Next()
		}
		return fiber.ErrUnauthorized
	}
}

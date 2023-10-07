package controllers

import (
	"chicchat/models"
	"chicchat/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandlerEmailLogin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		var user models.User
		err = db.Where("email = ?", req.Email).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Account not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		if !(user.Email == req.Email && utils.CheckPasswordHash(user.Password, req.Password)) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Incorrect email or password"})
		}

		accessToken, err := utils.CreateAccessToken(user.Email, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		refreshToken, err := utils.CreateRefreshToken(user.Email, db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}
		c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: refreshToken, HTTPOnly: true})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token": accessToken,
			"msg":          "Login success",
			"user":         user,
		})
	}
}
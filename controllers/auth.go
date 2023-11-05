package controllers

import (
	"chicchat/models"
	"chicchat/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func HandlerUserLogin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		var user models.User
		err = db.Where("username = ?", req.Username).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Account not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		if !(user.Username == req.Username && utils.CheckPasswordHash(user.Password, req.Password)) {
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

		respondUser := utils.RemoveUserSensitiveData(user)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token": accessToken,
			"msg":          "Login success",
			"user":         respondUser,
		})
	}
}

func RegisterUserByEmail(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		err := c.BodyParser(&user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		token, err := utils.GetBearerToken(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}
		userProfile, err := utils.GetUserProfileFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}
		if userProfile.ID != user.ID {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"msg": "Unauthorized"})
		}

		user.Provider = "email"
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}
		user.ProfilePicture = "./assets/image/default.png"

		err = CreateUser(db, &user)
		if err != nil {
			if err.Error() == "Email already exist" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
			}
			if err.Error() == "Username already exist" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		respondUser := utils.RemoveUserSensitiveData(user)
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"user": respondUser,
			"msg":  "Create User Successfully",
		})
	}
}

package controllers

import (
	"chicchat/models"
	"chicchat/utils"
	"errors"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	var existUser models.User
	err := db.Where("email = ? OR username = ? ", user.Email, user.Username).First(&existUser).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existUser.Email == user.Email {
			return errors.New("Email already exist")
		}
		if existUser.Username == user.Username {
			return errors.New("Username already exist")
		}
		return err
	}

	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User
		err := c.BodyParser(&user)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		if user.Password != "" {
			user.Password, err = utils.HashPassword(user.Password)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
			}
		}
		err = db.Model(&user).Where("id = ?", user.ID).Updates(&user).Error

		respondUser := utils.RemoveUserSensitiveData(user)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": respondUser,
			"msg":  "Update user successfully",
		})
	}
}

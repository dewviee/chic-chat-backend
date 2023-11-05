package controllers

import (
	"chicchat/models"
	"chicchat/utils"
	"errors"
	"os"
	"strconv"

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

func GetUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stringUserID := c.Params("id")
		userID, err := strconv.Atoi(stringUserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		var user models.User
		err = db.Where("id = ?", userID).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": err.Error()})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}
		return c.Status(200).JSON(fiber.Map{
			"user": utils.RemoveUserSensitiveData(user),
			"msg":  "Get user successfully"})
	}
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
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		respondUser := utils.RemoveUserSensitiveData(user)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": respondUser,
			"msg":  "Update user successfully",
		})
	}
}

func UploadProfilePicture(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := utils.GetBearerToken(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		userProfile, err := utils.GetUserProfileFromToken(token)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		file, err := c.FormFile("profile_picture")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
		}

		fileByte, err := utils.ConvertImageFileToBytes(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		filePath, err := utils.CreateImageFile(fileByte, userProfile)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		var updateData models.User
		updateData.ProfilePicture = filePath
		err = db.Model(&models.User{}).Where("id = ?", userProfile.ID).Updates(&updateData).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"msg":  "Create profile picture successfully",
			"user": utils.RemoveUserSensitiveData(updateData),
		})
	}
}

func GetProfilePicture(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filename := c.Params("filename")
		filePath := "./assets/image/" + filename

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File does not exist, send the default image
			defaultFilePath := "./assets/image/default.png"
			return c.SendFile(defaultFilePath)
		}

		// File exists, send the requested file
		return c.SendFile(filePath)
	}
}

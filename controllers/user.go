package controllers

import (
	"chicchat/models"
	"errors"

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

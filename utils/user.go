package utils

import (
	"chicchat/models"
	"time"
)

type respondUser struct {
	ID                 uint       `json:"id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	Phone              string     `json:"phone"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	ProfilePictureName string     `json:"profile_picture_name"`
}

func RemoveUserSensitiveData(user models.User) respondUser {
	var newUser respondUser
	newUser.ID = user.ID
	newUser.Username = user.Username
	newUser.Email = user.Email
	newUser.Phone = user.Phone
	newUser.CreatedAt = user.CreatedAt
	newUser.UpdatedAt = user.UpdatedAt
	newUser.ProfilePictureName = user.ProfilePictureName
	return newUser
}

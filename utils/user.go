package utils

import (
	"chicchat/models"
	"time"
)

type respondUser struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func RemoveUserSensitiveData(user models.User) respondUser {
	var newUser respondUser
	newUser.ID = user.ID
	newUser.Username = user.Username
	newUser.Email = user.Email
	newUser.CreatedAt = user.CreatedAt
	newUser.UpdatedAt = user.UpdatedAt
	return newUser
}

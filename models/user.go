package models

import "time"

type User struct {
	ID             uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Email          string     `json:"email" gorm:"unique;notnull"`
	Username       string     `json:"username" gorm:"unique;notnull"`
	Phone          string     `json:"phone"`
	Password       string     `json:"password"`
	CreatedAt      *time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Provider       string     `json:"provider" gorm:"notnull"`
	ProfilePicture string     `json:"profile_picture"`
}

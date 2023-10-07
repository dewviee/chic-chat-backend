package models

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Room struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Code string `json:"code" gorm:"unique;not null"`
}

type RoomWebSocket struct {
	ID        string
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mu        sync.Mutex
}

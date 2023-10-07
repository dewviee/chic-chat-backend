package models

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type RoomWebSocket struct {
	ID        string
	Clients   map[*websocket.Conn]bool
	Broadcast chan []byte
	Mu        sync.Mutex
}

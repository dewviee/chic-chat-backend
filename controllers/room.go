// controllers/room.go

package controllers

import (
	"chicchat/models"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type MessagePayload struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

func WebSocketHandler(rooms map[string]*models.RoomWebSocket) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		roomID := c.Params("id")
		room := rooms[roomID] // Get the room based on roomID

		if room == nil {
			// Handle the case when the room does not exist
			log.Printf("Room %s does not exist", roomID)
			c.Close()
			return
		}
		room.Mu.Lock()
		room.Clients[c] = true
		room.Mu.Unlock()

		log.Printf("Client connected to room %s", room.ID)

		defer func() {
			room.Mu.Lock()
			delete(room.Clients, c)
			room.Mu.Unlock()
			c.Close()
			log.Printf("Client disconnected from room %s", room.ID)
		}()

		for {
			_, byteMessage, err := c.ReadMessage()
			if err != nil {
				log.Println("Error when read message:", err)
				break
			}

			var message MessagePayload
			err = json.Unmarshal(byteMessage, &message)
			if err != nil {
				log.Println("JSON unmarshal error:", err)
				break
			}

			log.Printf("Received message from user '%s' in room '%s'", message.Message, room.ID)

			// Broadcast the message to all clients in the room
			room.Broadcast <- byteMessage
		}
	})
}

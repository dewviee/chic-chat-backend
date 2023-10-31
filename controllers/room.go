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

		if len(room.Clients) > 2 {
			log.Printf("Client disconnected from room %s", room.ID)
			return
		}

		room.Mu.Lock()
		room.Clients[c] = true
		room.Mu.Unlock()

		log.Printf("Client connected to room %s", room.ID)
		log.Printf("Client room %s have %d people", room.ID, len(room.Clients))

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

			log.Printf("Received message from user in room '%s'", room.ID)

			// Broadcast the message to all clients in the room
			room.Broadcast <- byteMessage
		}
	})
}

func GetRoomPeopleCount(rooms map[string]*models.RoomWebSocket) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		room := rooms[roomID]

		if room == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Room not found",
			})
		}

		if len(room.Clients) >= 2 {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Room is full",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Room Ready to connect",
		})
	}
}

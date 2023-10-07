// middlewares/room_websocket.go

package middlewares

import (
	"chicchat/models"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebSocketMiddleware(rooms map[string]*models.Room) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		if _, exists := rooms[roomID]; !exists {
			// Create a new room if it doesn't exist
			rooms[roomID] = &models.Room{
				ID:        roomID,
				Clients:   make(map[*websocket.Conn]bool),
				Broadcast: make(chan []byte),
			}

			// Start a goroutine to handle room broadcasting
			go func(room *models.Room) {
				for {
					message := <-room.Broadcast
					room.Mu.Lock()
					for client := range room.Clients {
						err := client.WriteMessage(websocket.TextMessage, message)
						if err != nil {
							log.Println("write:", err)
						}
					}
					room.Mu.Unlock()
				}
			}(rooms[roomID])
		}

		// IsWebSocketUpgrade returns true if the client requested an upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("room", rooms[roomID])
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

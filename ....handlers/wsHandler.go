package handlers

import (
	"archroid/archGap/config"
	"archroid/archGap/models"
	"archroid/archGap/services"
	"archroid/archGap/utils"
	"net/http"
	"sync"

	"github.com/charmbracelet/log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var userConnections = make(map[uint]*websocket.Conn)
var mutex = &sync.Mutex{}

type Input struct {
	Content     string `json:"content"`
	ReceiverID  uint   `json:"receiver_id"`
	MessageType string `json:"message_type"`
}

// Handle WebSocket connection for 1v1 chats
func HandleWebSocket(c echo.Context) error {
	// Get the user ID from the JWT token
	userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins (FOR DEVELOPMENT ONLY, restrict in production)
		},
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Lock and store the user's connection
	mutex.Lock()
	userConnections[userID] = conn
	mutex.Unlock()

	// Listen for messages from the user
	for {
		var input Input
		if err := conn.ReadJSON(&input); err != nil {
			log.Error("Read error:", err)
			log.Error("Error reading message:", err)
			break
		}

		// Parse and save message to the DB
		messageModel := models.Message{
			SenderID:    userID,
			Content:     input.Content,
			MessageType: input.MessageType,
		}

		// Save to DB
		err = services.SaveMessage(config.DB, messageModel)
		if err != nil {
			log.Error("Error saving message:", err)
		}

		// Send the message to the receiver's WebSocket (this is simplified; implement dynamic handling)
		receiverConn, ok := userConnections[messageModel.SenderID]
		if ok {
			err := receiverConn.WriteMessage(websocket.TextMessage, []byte(input.Content))
			if err != nil {
				log.Error("Error sending message:", err)
			}
		}
	}

	// Clean up connection on disconnect
	mutex.Lock()
	delete(userConnections, userID)
	mutex.Unlock()

	return nil
}

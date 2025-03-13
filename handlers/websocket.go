package handlers

import (
	"archroid/archGap/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var userID uint
var mutex = &sync.Mutex{}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections for simplicity
	},
}

type WebSocketMessage struct {
	Type   string `json:"type"`
	ChatID uint   `json:"chatID"`
	Content string `json:"content"`
}

type WebSocketConnection struct {
	Conn   *websocket.Conn
	UserID uint
}

var (
	userConnections   = make(map[uint]*websocket.Conn)
	chatSubscriptions = make(map[uint]map[uint]*websocket.Conn) // map of chatID -> map of userID -> conn
	mu                sync.RWMutex
)

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

	for {
		// Read the message from the WebSocket
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		if messageType == websocket.TextMessage {
			var msg WebSocketMessage
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Println("Error unmarshalling message:", err)
				continue
			}

			// Handle subscription to a chat
			if msg.Type == "subscribe" && msg.ChatID != 0 {
				subscribeUserToChat(userID, msg.ChatID, conn)
			} else if msg.Type == "unsubscribe" && msg.ChatID != 0 {
				unsubscribeUserFromChat(userID, msg.ChatID)
			} else if msg.Type == "message" && msg.Content != "" {
				// Handle sending a message to the chat
				sendMessageToChat(msg.ChatID, fmt.Sprintf("User %d: %s", userID, msg.Content))
			}
		}
	}

	// Clean up connection on disconnect
	mutex.Lock()
	delete(userConnections, userID)
	mutex.Unlock()

	return nil
}

// Subscribe user to a chat
func subscribeUserToChat(userID uint, chatID uint, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the chat already exists in subscriptions, if not create it
	if chatSubscriptions[chatID] == nil {
		chatSubscriptions[chatID] = make(map[uint]*websocket.Conn)
	}

	// Add the user connection to the chat subscriptions
	chatSubscriptions[chatID][userID] = conn
	fmt.Printf("User %s subscribed to chat %d\n", userID, chatID)
}

// Unsubscribe user from a chat
func unsubscribeUserFromChat(userID uint, chatID uint) {
	mutex.Lock()
	defer mutex.Unlock()

	// Remove the user from chat subscriptions
	delete(chatSubscriptions[chatID], userID)
	fmt.Printf("User %s unsubscribed from chat %d\n", userID, chatID)
}

// Broadcast message to all users in a chat
func sendMessageToChat(chatID uint, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Loop through all users subscribed to the chat and send them the message
	for userID, conn := range chatSubscriptions[chatID] {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error sending message to user %s: %v", userID, err)
		}
	}
}

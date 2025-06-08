package handlers

import (
	"archroid/archGap/db"
	"archroid/archGap/models"
	"archroid/archGap/utils"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var mutex = &sync.Mutex{}

type WebSocketMessage struct {
	Type        string `json:"type"`
	ChatID      uint   `json:"chatID"`
	Content     string `json:"content"`
	Status      string `json:"status"`
	MessageType string `json:"messageType"`
	Token       string `json:"token"`
}

type WebSocketConnection struct {
	Conn   *websocket.Conn
	UserID uint
}

var (
	userConnections   = make(map[uint]*websocket.Conn)
	chatSubscriptions = make(map[uint]map[uint]*websocket.Conn)
)

func HandleWebSocket(c echo.Context) error {
	// Get the user ID from the JWT token
	// userID, err := utils.ParseJWT(c.Request().Header.Get("Authorization"))
	// log.Println(c.Request().Header.Get("Authorization"))
	var userID uint

	// if err != nil {
	// 	return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
	// }

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

	for {
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

			if msg.Type == "message" && msg.Content != "" {
				// Handle sending a message to the chat
				sendMessageToChat(msg.ChatID, msg.Content, msg.MessageType, userID)
			} else if msg.Type == "getmessages" {
				// Handle fetching messages for a chat
				messages, err := db.GetMessagesinChat(msg.ChatID)
				if err != nil {
					log.Println("Error fetching messages:", err)
					continue
				}
				sendMessagestoChat(msg.ChatID, messages)

			} else if msg.Type == "verify" {
				userID, err = utils.ParseJWT(msg.Token)
				if err != nil {
					log.Println(err)
				} else {
					// Lock and store the user's connection
					mutex.Lock()
					userConnections[userID] = conn
					mutex.Unlock()

				}

			} else if msg.Type == "subscribe" && msg.ChatID != 0 && userID != 0 {
				subscribeUserToChat(userID, msg.ChatID, conn)
			} else if msg.Type == "unsubscribe" && msg.ChatID != 0 && userID != 0 {
				unsubscribeUserFromChat(userID, msg.ChatID)
			}
		}
	}

	// Clean up connection on disconnect
	mutex.Lock()
	delete(userConnections, userID)
	mutex.Unlock()

	return nil
}

func subscribeUserToChat(userID uint, chatID uint, conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the chat already exists in subscriptions, if not create it
	if chatSubscriptions[chatID] == nil {
		chatSubscriptions[chatID] = make(map[uint]*websocket.Conn)
	}

	chatSubscriptions[chatID][userID] = conn
	log.Printf("User %d subscribed to chat %d\n", userID, chatID)
}

func unsubscribeUserFromChat(userID uint, chatID uint) {
	mutex.Lock()
	defer mutex.Unlock()

	// Remove the user from chat subscriptions
	delete(chatSubscriptions[chatID], userID)
	log.Printf("User %d unsubscribed from chat %d\n", userID, chatID)
}

func sendMessageToChat(chatID uint, message string, messageType string, userid uint) {
	mutex.Lock()
	defer mutex.Unlock()

	// Loop through all users subscribed to the chat and send them the message
	for userID, conn := range chatSubscriptions[chatID] {
		msg := map[string]interface{}{
			"type":     "message",
			"senderID": userid,
			"text":     message,
		}
		msgBytes, _ := json.Marshal(msg)

		err := db.SendMessage(chatID, userID, message, messageType)
		if err != nil {
			log.Printf("Error saving message to database for chat %d: %v", chatID, err)
		}

		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("Error sending message to user %d: %v", userID, err)
		}
	}
}

func sendMessagestoChat(chatID uint, messages []models.Message) {
	mutex.Lock()
	defer mutex.Unlock()

	// Loop through all users subscribed to the chat and send them the messages
	for userID, conn := range chatSubscriptions[chatID] {
		msg := map[string]interface{}{
			"type":     "messages",
			"messages": messages,
		}
		msgBytes, _ := json.Marshal(msg)

		err := conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Printf("Error sending messages to user %d: %v", userID, err)
		}
	}
}

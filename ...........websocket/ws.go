package websocket

import (
	"net/http"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mu sync.Mutex

type Message struct {
	UserID  uint   `json:"user_id"`
	Content string `json:"content"`
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	clients[conn] = true
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Error("Read error:", err)
			delete(clients, conn)
			break
		}
		broadcast <- msg
	}
}

func BroadcastMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				log.Error("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

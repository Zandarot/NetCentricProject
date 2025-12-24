package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatUser struct {
	Username string
	Conn     *websocket.Conn
}

var onlineUsers = make(map[string]*ChatUser)

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	fmt.Println("💬 WebSocket Chat Server :9200")
	fmt.Println("   Private: @username message")
	log.Fatal(http.ListenAndServe(":9200", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	username := r.URL.Query().Get("user")
	if username == "" {
		sendJSON(conn, map[string]interface{}{"type": "error", "content": "Need ?user=YourName"})
		return
	}

	if _, exists := onlineUsers[username]; exists {
		sendJSON(conn, map[string]interface{}{"type": "error", "content": "Name taken"})
		return
	}

	onlineUsers[username] = &ChatUser{Username: username, Conn: conn}
	fmt.Printf("👤 %s joined (%d online)\n", username, len(onlineUsers))

	sendJSON(conn, map[string]interface{}{
		"type":    "system",
		"content": fmt.Sprintf("Welcome %s!", username),
	})

	// Tell others
	for _, user := range onlineUsers {
		if user.Username != username {
			sendJSON(user.Conn, map[string]interface{}{
				"type":    "system",
				"content": fmt.Sprintf("%s joined", username),
			})
		}
	}

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("👋 %s left\n", username)
			delete(onlineUsers, username)
			// Notify others
			for _, user := range onlineUsers {
				sendJSON(user.Conn, map[string]interface{}{
					"type":    "system",
					"content": fmt.Sprintf("%s left", username),
				})
			}
			return
		}

		content, _ := msg["content"].(string)
		if strings.HasPrefix(content, "@") {
			handlePrivate(username, content)
		} else {
			broadcastMessage(username, content)
		}
	}
}

func handlePrivate(sender, content string) {
	parts := strings.SplitN(strings.TrimSpace(content), " ", 2)
	if len(parts) < 2 {
		sendJSON(onlineUsers[sender].Conn, map[string]interface{}{
			"type":    "error",
			"content": "Use: @username message",
		})
		return
	}

	to := strings.TrimPrefix(parts[0], "@")
	message := parts[1]

	if to == sender {
		sendJSON(onlineUsers[sender].Conn, map[string]interface{}{
			"type":    "error",
			"content": "Can't PM yourself",
		})
		return
	}

	recipient, exists := onlineUsers[to]
	if !exists {
		sendJSON(onlineUsers[sender].Conn, map[string]interface{}{
			"type":    "error",
			"content": fmt.Sprintf("%s offline", to),
		})
		return
	}

	// Send to recipient WITH SENDER NAME
	sendJSON(recipient.Conn, map[string]interface{}{
		"type":    "private",
		"sender":  sender,
		"content": message,
	})

	// Confirm to sender
	sendJSON(onlineUsers[sender].Conn, map[string]interface{}{
		"type":    "private_sent",
		"to":      to,
		"content": message,
	})

	fmt.Printf("📩 %s → %s\n", sender, to)
}

func broadcastMessage(sender, content string) {
	fmt.Printf("📢 %s: %s\n", sender, content)

	for username, user := range onlineUsers {
		// Create NEW map for each recipient to avoid reference issues
		msg := map[string]interface{}{
			"type":    "message",
			"content": content,
		}

		// Set sender based on recipient
		if username == sender {
			msg["sender"] = "You"
		} else {
			msg["sender"] = sender
		}

		sendJSON(user.Conn, msg)
	}
}

func sendJSON(conn *websocket.Conn, data interface{}) {
	conn.WriteJSON(data)
}

package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func Chat() {
	// Get username from token
	username, err := getUsernameFromToken()
	if err != nil {
		fmt.Print("Enter your chat name: ")
		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if username == "" {
			username = "Guest"
		}
	}

	fmt.Printf("🔌 Connecting as %s...\n", username)
	fmt.Println("💬 Commands:")
	fmt.Println("  @username message  - Send private message")
	fmt.Println("  normal message     - Send to everyone")
	fmt.Println("  /exit              - Quit chat")

	// Connect to WebSocket with username parameter
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:9200/ws?user="+username, nil)
	if err != nil {
		fmt.Printf(" Failed to connect: %v\n", err)
		fmt.Println("   Make sure WebSocket server is running: go run cmd/ws_server/main.go")
		return
	}
	defer conn.Close()

	fmt.Println(" Connected to chat!")

	// Receive messages in background
	go func() {
		for {
			var msg map[string]interface{}
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Println("\n🔌 Disconnected")
				os.Exit(0)
			}

			msgType, _ := msg["type"].(string)
			sender, _ := msg["sender"].(string)
			content, _ := msg["content"].(string)
			to, _ := msg["to"].(string)

			switch msgType {
			case "system":
				fmt.Printf("\033[90m %s\033[0m\n", content)
			case "message":
				fmt.Printf("\033[92m[%s] %s\033[0m\n", sender, content)
			case "private":
				fmt.Printf("\033[95m [PM from %s] %s\033[0m\n", sender, content)
			case "private_sent":
				fmt.Printf("\033[96m [PM to %s] %s\033[0m\n", to, content)
			case "error":
				fmt.Printf("\033[91m %s\033[0m\n", content)
			default:
				fmt.Printf("\033[93m[Unknown message type: %s] %s\033[0m\n", msgType, content)
			}
		}
	}()

	// Send messages
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("\n💭 Type your message: ")

	for scanner.Scan() {
		text := scanner.Text()

		if text == "/exit" || text == "/quit" {
			fmt.Println(" Goodbye!")
			break
		}

		if text == "" {
			fmt.Print(" Type your message: ")
			continue
		}

		// Send message to server
		msg := map[string]string{
			"content": text,
		}

		if err := conn.WriteJSON(msg); err != nil {
			fmt.Printf(" Send failed: %v\n", err)
			break
		}

		fmt.Print(" Type your message: ")
	}
}

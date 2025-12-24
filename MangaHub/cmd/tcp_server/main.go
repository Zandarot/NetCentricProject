package main

import (
	"MangaHub/internal/tcp"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var hub = tcp.NewHub()

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal("Unable to start TCP server", err)
	}
	fmt.Println("TCP is listening on port 9000") // FIXED: Was "9090"
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error: ", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	fmt.Println("Client connected : ", addr)

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected : ", addr)
			return
		}
		var msg tcp.Message
		err = json.Unmarshal(buf[:n], &msg)
		if err != nil {
			continue
		}
		if msg.Type == "register" {
			hub.Register(msg.UserID, conn)
			fmt.Println("Registered TCP client for user:", msg.UserID)
		}
		if msg.Type == "progress_sync" {
			data, _ := json.Marshal(msg)
			hub.Broadcast(msg.UserID, data)
		}
	}
}

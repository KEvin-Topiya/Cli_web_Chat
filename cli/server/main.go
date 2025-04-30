package main

import (
	"bufio"
	"fmt"
	"net"
)

var clients = make(map[net.Conn]bool)
var messages = make(chan string)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("✅ Server started on :8080")

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("❌ Error accepting:", err)
			continue
		}
		clients[conn] = true
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			delete(clients, conn)
			conn.Close()
			break
		}
		messages <- fmt.Sprintf("%s", message)
	}
}

func broadcast() {
	for {
		msg := <-messages
		for client := range clients {
			_, err := fmt.Fprint(client, msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

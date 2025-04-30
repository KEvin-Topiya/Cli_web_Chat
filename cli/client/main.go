package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Ask for username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1] // remove newline

	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected to chat server as", username)
	fmt.Println("👉 Start chatting below\n")

	// Receive messages
	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Disconnected from server")
				os.Exit(0)
			}
			fmt.Print(message)
		}
	}()

	// Send messages
	for {
		text, _ := reader.ReadString('\n')
		fullMessage := fmt.Sprintf("[%s]: %s", username, text)
		_, err := conn.Write([]byte(fullMessage))
		if err != nil {
			fmt.Println("❌ Failed to send message")
			break
		}
	}
}

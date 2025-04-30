package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string   `json:"type"`
	Username string   `json:"username,omitempty"`
	Text     string   `json:"text,omitempty"`
	NewName  string   `json:"newName,omitempty"`
	OldName  string   `json:"oldName,omitempty"`
	Users    []string `json:"users,omitempty"`
}

func main() {
	fmt.Print("Enter your username: ")
	var username string
	fmt.Scanln(&username)

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	conn.WriteJSON(Message{Type: "join", Username: username})

	go func() {
		for {
			var m Message
			if err := conn.ReadJSON(&m); err != nil {
				return
			}

			switch m.Type {
			case "chat":
				if m.Username == username {
					color.New(color.FgGreen).Printf("You: %s\n", m.Text)
				} else {
					color.New(color.FgCyan).Printf("%s: %s\n", m.Username, m.Text)
				}
			case "error":
				color.New(color.FgRed).Println("Error:", m.Text)
			case "rename_success":
				username = m.NewName
				color.Yellow("Name changed to %s", username)
			case "user_list":
				color.Magenta("👥 Online: %v\n", m.Users)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/rename" {
			fmt.Print("New username: ")
			scanner.Scan()
			newName := scanner.Text()
			conn.WriteJSON(Message{Type: "rename", OldName: username, NewName: newName})
		} else {
			conn.WriteJSON(Message{Type: "chat", Username: username, Text: text})
		}
	}
}

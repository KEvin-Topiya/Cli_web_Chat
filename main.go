package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Text     string `json:"text,omitempty"`
	OldName  string `json:"oldName,omitempty"`
	NewName  string `json:"newName,omitempty"`
}

var (
	clients   = make(map[*websocket.Conn]string)
	usernames = make(map[string]bool)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	var currentName string

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, ws)
			if currentName != "" {
				delete(usernames, currentName)
			}
			break
		}

		switch msg.Type {
		case "join":
			if usernames[msg.Username] {
				ws.WriteJSON(map[string]string{
					"type":    "error",
					"message": "Username already taken!",
				})
				continue
			}
			currentName = msg.Username
			clients[ws] = msg.Username
			usernames[msg.Username] = true

		case "chat":
			broadcastMessage := Message{
				Type:     "chat",
				Username: msg.Username,
				Text:     msg.Text,
			}
			broadcastJSON(broadcastMessage)

		case "rename":
			if usernames[msg.NewName] {
				ws.WriteJSON(map[string]string{
					"type":    "error",
					"message": "That name is already taken.",
				})
			} else {
				delete(usernames, msg.OldName)
				usernames[msg.NewName] = true
				clients[ws] = msg.NewName
				currentName = msg.NewName

				ws.WriteJSON(map[string]string{
					"type":    "rename_success",
					"newName": msg.NewName,
				})
			}
		}
	}
}

func broadcastJSON(msg Message) {
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Println("Broadcast error:", err)
			client.Close()
			delete(usernames, clients[client])
			delete(clients, client)
		}
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	log.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}

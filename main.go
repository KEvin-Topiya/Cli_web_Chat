package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

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

var (
	clients      = make(map[*websocket.Conn]string)
	clientsMutex = sync.Mutex{}
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	http.HandleFunc("/ws", handleConnections)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade:", err)
		return
	}
	defer conn.Close()

	var username string

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue
		}

		switch m.Type {
		case "join":
			clientsMutex.Lock()
			if usernameTaken(m.Username) {
				conn.WriteJSON(Message{Type: "error", Text: "Username already taken"})
				clientsMutex.Unlock()
				continue
			}
			username = m.Username
			clients[conn] = username
			clientsMutex.Unlock()

			broadcastUserList()

		case "chat":
			clientsMutex.Lock()
			broadcast(Message{Type: "chat", Username: m.Username, Text: m.Text})
			clientsMutex.Unlock()

		case "rename":
			clientsMutex.Lock()
			if usernameTaken(m.NewName) {
				conn.WriteJSON(Message{Type: "error", Text: "Username already taken"})
			} else {
				clients[conn] = m.NewName
				username = m.NewName
				conn.WriteJSON(Message{Type: "rename_success", NewName: m.NewName})
				broadcastUserList()
			}
			clientsMutex.Unlock()
		}
	}

	clientsMutex.Lock()
	delete(clients, conn)
	clientsMutex.Unlock()
	broadcastUserList()
}

func broadcast(msg Message) {
	for conn := range clients {
		conn.WriteJSON(msg)
	}
}

func broadcastUserList() {
	users := []string{}
	for _, name := range clients {
		users = append(users, name)
	}
	broadcast(Message{Type: "user_list", Users: users})
}

func usernameTaken(name string) bool {
	for _, n := range clients {
		if n == name {
			return true
		}
	}
	return false
}

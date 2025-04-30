
# 💬 CLI & Web Chat Application (Go WebSockets)

A real-time chat application supporting both **CLI** and **Web** interfaces using a single WebSocket server in Go. Users can chat with each other, change their username, and view active users — all in real time.

---

## 📁 Project Structure

```
cli_web_chat/
├── cli/               # CLI chat client
│   └── main.go
├── static/            # Web frontend (HTML, CSS, JS)
│   └── index.html
├── main.go            # Go WebSocket server
└── README.md
```

---

## 🚀 Features

- ✅ Real-time communication via WebSocket
- ✅ Web and CLI clients using the same server
- ✅ Unique username enforcement
- ✅ Username editing support (`/rename` in CLI, clickable in web)
- ✅ Emoji and color support in CLI
- ✅ Live user list sync across clients

---

## 🧑‍💻 Requirements

- Go 1.17 or newer
- A modern web browser
- Terminal (for CLI usage)
- Dependencies:
  - `github.com/gorilla/websocket`
  - `github.com/fatih/color`

---

## 🛠️ Setup & Usage

### 1. Clone the Repository

```bash
git clone https://github.com/KEvin-Topiya/Cli_web_Chat.git
cd Cli_web_Chat
```

### 2. Install Dependencies

```bash
go get github.com/gorilla/websocket
go get github.com/fatih/color
```

### 3. Run the WebSocket Server

```bash
go run main.go
```

### 4. Open Web Chat in Browser

Navigate to:

```
http://localhost:8080
```

### 5. Run CLI Chat in New Terminal

```bash
cd cli
go run main.go
```

---

## 🧪 CLI Commands

- `/rename` – Change your username

You can also type any message and press `Enter` to send.

---

## 📺 Preview

- 🌐 Web: Clean, responsive chat interface with light theme  
- 🖥️ CLI: Colored usernames, emoji support, real-time chat

---

## 🔐 Notes

- You can start either **CLI or Web first** – they work together in any order.
- Usernames must be unique.
- Live user list is synchronized across all clients.

---

## 📌 Future Enhancements

- Persistent chat logs
- Authentication
- Docker support
- Mobile UI optimization
- Emoji picker in Web UI

---

## 👨‍💻 Author

Made with ❤️ by Kevin (2025)  
Project: `cli_web_chat`

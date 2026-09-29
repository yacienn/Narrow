package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents one connected player.
type Client struct {
	conn *websocket.Conn
	id   int
}

// Game manages all connected players.
type Game struct {
	clients map[int]*Client
	nextID  int
	mu      sync.Mutex
}

var game = Game{
	clients: make(map[int]*Client),
	nextID:  1,
}

// Messages sent between Godot and Go.
type Message struct {
	Type string `json:"type"`
}

// Message sent when a player connects.
type WelcomeMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

func gameHandler(w http.ResponseWriter, r *http.Request) {

	// Upgrade HTTP connection to WebSocket.
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket upgrade error:", err)
		return
	}

	// Give the player a unique ID.
	game.mu.Lock()

	playerID := game.nextID
	game.nextID++

	client := &Client{
		conn: conn,
		id:   playerID,
	}

	game.clients[playerID] = client

	game.mu.Unlock()

	fmt.Println("Player connected:", playerID)

	// Tell Godot its player ID.
	welcome := WelcomeMessage{
		Type:     "welcome",
		PlayerID: playerID,
	}

	data, err := json.Marshal(welcome)
	if err != nil {
		fmt.Println("json error:", err)
		conn.Close()
		return
	}

	err = conn.WriteMessage(
		websocket.TextMessage,
		data,
	)

	if err != nil {
		fmt.Println("welcome message error:", err)
		removeClient(playerID)
		return
	}

	// Keep listening for messages from this player.
	for {

		messageType, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Player disconnected:", playerID)

			removeClient(playerID)

			return
		}

		fmt.Printf(
			"Player %d sent: %s\n",
			playerID,
			string(message),
		)

		// Temporary echo.
		err = conn.WriteMessage(
			messageType,
			message,
		)

		if err != nil {
			fmt.Println("write error:", err)

			removeClient(playerID)

			return
		}
	}
}

func removeClient(playerID int) {

	game.mu.Lock()
	defer game.mu.Unlock()

	client, exists := game.clients[playerID]

	if !exists {
		return
	}

	client.conn.Close()

	delete(game.clients, playerID)

	fmt.Println("Removed player:", playerID)
}

func main() {

	http.HandleFunc("/game", gameHandler)

	fmt.Println("Server listening on :8080")

	err := http.ListenAndServe(
		"0.0.0.0:8080",
		nil,
	)

	if err != nil {
		fmt.Println("server error:", err)
	}
}

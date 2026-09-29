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

type Client struct {
	conn *websocket.Conn
	id   int
}

type Game struct {
	clients map[int]*Client
	nextID  int
	mu      sync.Mutex
}

var game = Game{
	clients: make(map[int]*Client),
	nextID:  1,
}

type Message struct {
	Type string `json:"type"`
}

type WelcomeMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

type PlayerJoinedMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

func gameHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("websocket upgrade error:", err)
		return
	}

	// Create a new player
	game.mu.Lock()

	playerID := game.nextID
	game.nextID++

	client := &Client{
		conn: conn,
		id:   playerID,
	}

	game.clients[playerID] = client

	// Save a list of players that were already connected.
	// We need this before unlocking.
	existingPlayers := make([]int, 0, len(game.clients))

	for id := range game.clients {
		if id != playerID {
			existingPlayers = append(existingPlayers, id)
		}
	}

	game.mu.Unlock()

	fmt.Println("Player connected:", playerID)

	// Tell the new player its ID.
	welcome := WelcomeMessage{
		Type:     "welcome",
		PlayerID: playerID,
	}

	sendJSON(conn, welcome)

	// Tell the new player about existing players.
	for _, existingID := range existingPlayers {

		message := PlayerJoinedMessage{
			Type:     "player_joined",
			PlayerID: existingID,
		}

		sendJSON(conn, message)
	}

	// Tell existing players that this player joined.
	broadcastToOthers(playerID, PlayerJoinedMessage{
		Type:     "player_joined",
		PlayerID: playerID,
	})

	// Listen for messages from this player.
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
		err = conn.WriteMessage(messageType, message)

		if err != nil {
			fmt.Println("write error:", err)

			removeClient(playerID)

			return
		}
	}
}

func sendJSON(conn *websocket.Conn, message any) {

	data, err := json.Marshal(message)

	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	err = conn.WriteMessage(
		websocket.TextMessage,
		data,
	)

	if err != nil {
		fmt.Println("send error:", err)
	}
}

func broadcastToOthers(senderID int, message any) {

	data, err := json.Marshal(message)

	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	game.mu.Lock()
	defer game.mu.Unlock()

	for id, client := range game.clients {

		if id == senderID {
			continue
		}

		err := client.conn.WriteMessage(
			websocket.TextMessage,
			data,
		)

		if err != nil {
			fmt.Println(
				"broadcast error to player",
				id,
				":",
				err,
			)
		}
	}
}

func removeClient(playerID int) {

	game.mu.Lock()

	client, exists := game.clients[playerID]

	if !exists {
		game.mu.Unlock()
		return
	}

	delete(game.clients, playerID)

	game.mu.Unlock()

	client.conn.Close()

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

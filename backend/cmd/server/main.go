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

	// gorilla/websocket allows only one writer at a time per connection.
	writeMu sync.Mutex
}

// send writes a text message safely (one writer at a time).
func (c *Client) send(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteMessage(websocket.TextMessage, data)
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

// Anything the client sends us.
type IncomingMessage struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

type WelcomeMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

type PlayerJoinedMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

type PlayerLeftMessage struct {
	Type     string `json:"type"`
	PlayerID int    `json:"player_id"`
}

type PlayerMovedMessage struct {
	Type     string  `json:"type"`
	PlayerID int     `json:"player_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
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
	existingPlayers := make([]int, 0, len(game.clients))

	for id := range game.clients {
		if id != playerID {
			existingPlayers = append(existingPlayers, id)
		}
	}

	game.mu.Unlock()

	fmt.Println("Player connected:", playerID)

	// Tell the new player its ID.
	sendJSON(client, WelcomeMessage{
		Type:     "welcome",
		PlayerID: playerID,
	})

	// Tell the new player about existing players.
	for _, existingID := range existingPlayers {
		sendJSON(client, PlayerJoinedMessage{
			Type:     "player_joined",
			PlayerID: existingID,
		})
	}

	// Tell existing players that this player joined.
	broadcastToOthers(playerID, PlayerJoinedMessage{
		Type:     "player_joined",
		PlayerID: playerID,
	})

	// Listen for messages from this player.
	for {

		_, raw, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Player disconnected:", playerID)
			removeClient(playerID)
			return
		}

		var msg IncomingMessage

		if err := json.Unmarshal(raw, &msg); err != nil {
			fmt.Println("invalid JSON from player", playerID, ":", err)
			continue
		}

		switch msg.Type {

		case "move":
			// Relay this player's position to everyone else.
			broadcastToOthers(playerID, PlayerMovedMessage{
				Type:     "player_moved",
				PlayerID: playerID,
				X:        msg.X,
				Y:        msg.Y,
			})

		default:
			fmt.Printf("Player %d sent unknown type: %s\n", playerID, msg.Type)
		}
	}
}

func sendJSON(client *Client, message any) {

	data, err := json.Marshal(message)

	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	if err := client.send(data); err != nil {
		fmt.Println("send error:", err)
	}
}

func broadcastToOthers(senderID int, message any) {

	data, err := json.Marshal(message)

	if err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	// Copy the recipients so we don't hold the lock while writing.
	game.mu.Lock()
	recipients := make([]*Client, 0, len(game.clients))
	for id, client := range game.clients {
		if id != senderID {
			recipients = append(recipients, client)
		}
	}
	game.mu.Unlock()

	for _, client := range recipients {
		if err := client.send(data); err != nil {
			fmt.Println("broadcast error to player", client.id, ":", err)
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

	// Tell everyone else this player left.
	broadcastToOthers(playerID, PlayerLeftMessage{
		Type:     "player_left",
		PlayerID: playerID,
	})
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

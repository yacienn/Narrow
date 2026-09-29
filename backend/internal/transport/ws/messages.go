package ws

import (
	"fmt"
	"encoding/json"

	"github.com/yacienn/Game/internal/domain"
)

// --- Client -> server ---

type incomingMessage struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

// --- Server -> client (wire format; unchanged from the original) ---

type welcomeMessage struct {
	Type     string          `json:"type"`
	PlayerID domain.PlayerID `json:"player_id"`
}

type playerJoinedMessage struct {
	Type     string          `json:"type"`
	PlayerID domain.PlayerID `json:"player_id"`
}

type playerLeftMessage struct {
	Type     string          `json:"type"`
	PlayerID domain.PlayerID `json:"player_id"`
}

type playerMovedMessage struct {
	Type     string          `json:"type"`
	PlayerID domain.PlayerID `json:"player_id"`
	X        float64         `json:"x"`
	Y        float64         `json:"y"`
}

// encode maps a domain event to its JSON wire representation.
func encode(ev domain.Event) ([]byte, error) {
	switch e := ev.(type) {
	case domain.Welcomed:
		return json.Marshal(welcomeMessage{Type: "welcome", PlayerID: e.PlayerID})
	case domain.PlayerJoined:
		return json.Marshal(playerJoinedMessage{Type: "player_joined", PlayerID: e.PlayerID})
	case domain.PlayerLeft:
		return json.Marshal(playerLeftMessage{Type: "player_left", PlayerID: e.PlayerID})
	case domain.PlayerMoved:
		return json.Marshal(playerMovedMessage{
			Type:     "player_moved",
			PlayerID: e.PlayerID,
			X:        e.Position.X,
			Y:        e.Position.Y,
		})
	default:
		return nil, fmt.Errorf("unsupported event %T", ev)
	}
}

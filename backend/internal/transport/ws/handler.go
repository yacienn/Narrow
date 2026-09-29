// Package ws is the WebSocket delivery layer. It translates between JSON
// frames and calls on the GameService.
package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/yacienn/Game/internal/domain"
	"github.com/yacienn/Game/internal/usecase"
)

type Handler struct {
	upgrader websocket.Upgrader
	game     *usecase.GameService
	hub      *Hub
}

func NewHandler(game *usecase.GameService, hub *Hub) *Handler {
	return &Handler{
		upgrader: websocket.Upgrader{
			// Fine for development; restrict origins in production.
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		game: game,
		hub:  hub,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade", "err", err)
		return
	}

	id := h.game.NextPlayerID()
	h.hub.Register(id, conn)
	slog.Info("player connected", "player", id)

	defer func() {
		h.hub.Unregister(id)
		h.game.Leave(id)
		slog.Info("player disconnected", "player", id)
	}()

	h.game.Join(id)
	h.readLoop(id, conn)
}

func (h *Handler) readLoop(id domain.PlayerID, conn *websocket.Conn) {
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg incomingMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			slog.Warn("invalid JSON", "player", id, "err", err)
			continue
		}

		switch msg.Type {
		case "move":
			h.game.Move(id, domain.Position{X: msg.X, Y: msg.Y})
		default:
			slog.Warn("unknown message type", "player", id, "type", msg.Type)
		}
	}
}

package ws

import (
	"log/slog"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/yacienn/Game/internal/domain"
	"github.com/yacienn/Game/internal/usecase"
)

var _ usecase.Notifier = (*Hub)(nil)

// client wraps a connection. gorilla/websocket allows only one writer at a
// time per connection, hence the mutex.
type client struct {
	id      domain.PlayerID
	conn    *websocket.Conn
	writeMu sync.Mutex
}

func (c *client) send(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Hub tracks live connections and implements usecase.Notifier.
type Hub struct {
	mu      sync.RWMutex
	clients map[domain.PlayerID]*client
}

func NewHub() *Hub {
	return &Hub{clients: make(map[domain.PlayerID]*client)}
}

func (h *Hub) Register(id domain.PlayerID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[id] = &client{id: id, conn: conn}
}

func (h *Hub) Unregister(id domain.PlayerID) {
	h.mu.Lock()
	c, ok := h.clients[id]
	delete(h.clients, id)
	h.mu.Unlock()

	if ok {
		c.conn.Close()
	}
}

func (h *Hub) SendTo(id domain.PlayerID, ev domain.Event) {
	data, err := encode(ev)
	if err != nil {
		slog.Error("encode event", "err", err)
		return
	}

	h.mu.RLock()
	c, ok := h.clients[id]
	h.mu.RUnlock()
	if !ok {
		return
	}

	if err := c.send(data); err != nil {
		slog.Warn("send failed", "player", id, "err", err)
	}
}

func (h *Hub) BroadcastExcept(id domain.PlayerID, ev domain.Event) {
	data, err := encode(ev)
	if err != nil {
		slog.Error("encode event", "err", err)
		return
	}

	// Copy recipients so we don't hold the lock while writing.
	h.mu.RLock()
	recipients := make([]*client, 0, len(h.clients))
	for cid, c := range h.clients {
		if cid != id {
			recipients = append(recipients, c)
		}
	}
	h.mu.RUnlock()

	for _, c := range recipients {
		if err := c.send(data); err != nil {
			slog.Warn("broadcast failed", "player", c.id, "err", err)
		}
	}
}

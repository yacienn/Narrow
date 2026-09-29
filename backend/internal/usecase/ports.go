// Package usecase contains the application logic. It depends only on the
// domain package and on the interfaces (ports) declared here.
package usecase

import "github.com/yacienn/Game/internal/domain"

// PlayerRepository is implemented by the storage layer.
type PlayerRepository interface {
	NextID() domain.PlayerID
	Add(p domain.Player)
	// Remove reports whether the player existed.
	Remove(id domain.PlayerID) bool
	// UpdatePosition reports whether the player exists.
	UpdatePosition(id domain.PlayerID, pos domain.Position) bool
	ListExcept(id domain.PlayerID) []domain.Player
}

// Notifier is implemented by the delivery layer (WebSocket, etc.).
type Notifier interface {
	SendTo(id domain.PlayerID, ev domain.Event)
	BroadcastExcept(id domain.PlayerID, ev domain.Event)
}

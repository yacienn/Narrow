// Package memory is an in-memory implementation of usecase.PlayerRepository.
// Swap it for Redis/Postgres later without touching the use cases.
package memory

import (
	"sync"

	"github.com/yacienn/Game/internal/domain"
	"github.com/yacienn/Game/internal/usecase"
)

var _ usecase.PlayerRepository = (*PlayerRepository)(nil)

type PlayerRepository struct {
	mu      sync.RWMutex
	nextID  domain.PlayerID
	players map[domain.PlayerID]domain.Player
}

func NewPlayerRepository() *PlayerRepository {
	return &PlayerRepository{
		nextID:  1,
		players: make(map[domain.PlayerID]domain.Player),
	}
}

func (r *PlayerRepository) NextID() domain.PlayerID {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID
	r.nextID++
	return id
}

func (r *PlayerRepository) Add(p domain.Player) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.players[p.ID] = p
}

func (r *PlayerRepository) Remove(id domain.PlayerID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.players[id]; !ok {
		return false
	}
	delete(r.players, id)
	return true
}

func (r *PlayerRepository) UpdatePosition(id domain.PlayerID, pos domain.Position) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.players[id]
	if !ok {
		return false
	}
	p.Position = pos
	r.players[id] = p
	return true
}

func (r *PlayerRepository) ListExcept(id domain.PlayerID) []domain.Player {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]domain.Player, 0, len(r.players))
	for pid, p := range r.players {
		if pid != id {
			out = append(out, p)
		}
	}
	return out
}

package usecase

import "github.com/yacienn/Game/internal/domain"

type GameService struct {
	players  PlayerRepository
	notifier Notifier
}

func NewGameService(players PlayerRepository, notifier Notifier) *GameService {
	return &GameService{players: players, notifier: notifier}
}

// NextPlayerID reserves an ID so the transport can register the connection
// before Join starts sending messages to it.
func (s *GameService) NextPlayerID() domain.PlayerID {
	return s.players.NextID()
}

// Join adds the player, tells them who they are and who is already here,
// and tells everyone else that they arrived.
func (s *GameService) Join(id domain.PlayerID) {
	s.players.Add(domain.Player{ID: id})
	existing := s.players.ListExcept(id)

	s.notifier.SendTo(id, domain.Welcomed{PlayerID: id})
	for _, p := range existing {
		s.notifier.SendTo(id, domain.PlayerJoined{PlayerID: p.ID})
	}
	s.notifier.BroadcastExcept(id, domain.PlayerJoined{PlayerID: id})
}

// Move stores the new position and relays it to everyone else.
func (s *GameService) Move(id domain.PlayerID, pos domain.Position) {
	if !s.players.UpdatePosition(id, pos) {
		return
	}
	s.notifier.BroadcastExcept(id, domain.PlayerMoved{PlayerID: id, Position: pos})
}

// Leave removes the player and tells everyone else.
func (s *GameService) Leave(id domain.PlayerID) {
	if !s.players.Remove(id) {
		return
	}
	s.notifier.BroadcastExcept(id, domain.PlayerLeft{PlayerID: id})
}

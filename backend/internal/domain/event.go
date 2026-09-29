package domain

// Event is something that happened in the game that other players must hear
// about. The use-case layer emits events; the transport layer decides how to
// serialize them (JSON today, anything tomorrow).
type Event interface{ isEvent() }

// Welcomed is sent to a player right after joining, telling them their ID.
type Welcomed struct{ PlayerID PlayerID }

type PlayerJoined struct{ PlayerID PlayerID }

type PlayerMoved struct {
	PlayerID PlayerID
	Position Position
}

type PlayerLeft struct{ PlayerID PlayerID }

func (Welcomed) isEvent()     {}
func (PlayerJoined) isEvent() {}
func (PlayerMoved) isEvent()  {}
func (PlayerLeft) isEvent()   {}

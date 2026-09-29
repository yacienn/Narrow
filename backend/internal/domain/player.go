// Package domain holds the core game entities. It imports nothing from the
// rest of the project (and nothing from any framework).
package domain

type PlayerID int

type Position struct {
	X float64
	Y float64
}

type Player struct {
	ID       PlayerID
	Position Position
}

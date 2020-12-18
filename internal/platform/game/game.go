package game

import "minesweeper/internal"

type Factory = func() Game

// A Game represents the internal state of a Minesweeper game
type Game interface {
	// IsStarted reports if the game has already started (a first call to Sync
	// is required)
	IsStarted() bool

	// Finish marks this game as finished
	Finish()

	// IsFinished queries if this game finished
	IsFinished() bool

	// Board gets a reference to the current Board, makes no guarantee about
	// internal state of the Game (domain "leaks", maybe fix eventually)
	Board() (internal.Board, error)

	// Sync this Game's Board with the given one.
	// This is a mutating, blocking operation
	Sync(internal.Board) error
}

func CanPlay(game Game) bool {
	return game.IsStarted() && !game.IsFinished()
}

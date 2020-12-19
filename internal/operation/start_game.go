package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/random"
)

type StartGame struct {
	Rand random.Intn
}

func (sg StartGame) StartGame(game game.Game, rows, columns, bombs int) error {
	if game.IsStarted() {
		return internal.NewInvalidOperation("game already started")
	}

	// TODO: Move this validation to Game
	if rows <= 0 || columns <= 0 || bombs <= 0 {
		return internal.NewInvalidOperation(fmt.Sprintf("creating game with %dx%d board and %d bombs", rows, columns, bombs))
	}

	board := internal.NewBoard(sg.Rand, rows, columns, bombs)
	return game.Sync(board)
}

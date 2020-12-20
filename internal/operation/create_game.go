package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/random"
)

type CreateGame struct {
	Holder game.Holder
	Rand   random.Intn
}

func (cg CreateGame) CreateGame(rows, columns, bombs int) (game.ID, internal.Board, error) {
	if rows <= 0 || columns <= 0 || bombs <= 0 {
		return 0, internal.Board{}, invalidBoardError(rows, columns, bombs)
	}

	board := internal.NewBoard(cg.Rand, rows, columns, bombs)
	g := game.NewInMemory(&board)

	id, err := cg.Holder.Insert(g)

	if err != nil {
		return 0, internal.Board{}, err
	}

	return id, board, nil
}

func invalidBoardError(rows, columns, bombs int) error {
	return internal.NewInvalidOperation(fmt.Sprintf("creating game with %dx%d board and %d bombs", rows, columns, bombs))
}

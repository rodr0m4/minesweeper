package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type Mark struct{}

func (m Mark) Mark(game game.Game, row, column int, mark internal.TileMark) error {
	board, err := game.Board()

	if err != nil {
		return err
	}

	position, err := board.Position(row, column)

	if err != nil {
		return err
	}

	tile := board.Find(position)

	if err := tile.Mark(mark); err != nil {
		return err
	}

	return game.Sync(board)
}

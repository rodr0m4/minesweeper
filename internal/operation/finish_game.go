package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type FinishGame struct{}

func (fg FinishGame) Finish(game game.Game) error {
	board, err := game.Board()

	if err != nil {
		return err
	}

	board.Traverse(func(tile *internal.Tile, position internal.Position) {
		tile.Modify(internal.Reveal())
	})

	return game.Sync(board)
}

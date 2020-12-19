package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type Tap struct {
	GameFinisher GameFinisher
	TileRevealer TileRevealer
}

type GameFinisher interface {
	Finish(game.Game) error
}

type TileRevealer interface {
	RevealAdjacent(game game.Game, board *internal.Board, position internal.Position) error
}

func (t Tap) Tap(game game.Game, row, column int) (result internal.TapResult, err error) {
	board, err := game.Board()

	if err != nil {
		return
	}

	position, err := board.Position(row, column)

	if err != nil {
		return
	}

	result, err = board.Find(position).Tap()

	if err != nil {
		return
	}

	switch result {
	case internal.ExplosionResult:
		err = t.GameFinisher.Finish(game)
	case internal.NothingResult:
		err = t.TileRevealer.RevealAdjacent(game, &board, position)
	default:
		panic(fmt.Sprintf("unrecheable code, invalid tap result %d", result))
	}

	return
}

package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type Tap struct {
	GameFinisher GameFinisher
}

type GameFinisher interface {
	Finish(game.Game) error
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

	if result == internal.ExplosionResult {
		err = t.GameFinisher.Finish(game)
	}

	return
}

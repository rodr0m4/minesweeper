package operation

import (
	"minesweeper/internal/platform/game"
)

type DeleteGame struct {
	Holder      game.Holder
	BoardDrawer BoardDrawer
}

func (dg DeleteGame) DeleteGame(id game.ID) (ShowedGame, error) {
	g, err := dg.Holder.Get(id)

	if err != nil {
		return ShowedGame{}, err
	}

	board, err := g.Board()

	if err != nil {
		return ShowedGame{}, err
	}

	if err := dg.Holder.Delete(id); err != nil {
		return ShowedGame{}, err
	}

	return dg.BoardDrawer.DrawBoard(board), nil
}

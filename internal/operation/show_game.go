package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type ShowGame struct {
	BoardDrawer boardDrawer
}

type boardDrawer interface {
	DrawBoardIntoStringArray(board internal.Board, revealEverything bool) []string
}

type ShowedGame struct {
	Lines []string `json:"lines"`
}

func (sg ShowGame) ShowGame(game game.Game) (ShowedGame, error) {
	board, err := game.Board()

	if err != nil {
		return ShowedGame{}, err
	}

	return ShowedGame{
		Lines: sg.BoardDrawer.DrawBoardIntoStringArray(board, true),
	}, nil
}

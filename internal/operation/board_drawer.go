package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type ShowedGame struct {
	Lines []string `json:"lines"`
}

type BoardDrawer interface {
	DrawBoard(internal.Board) ShowedGame
}

type DefaultBoardDrawer struct {
	RevealEverything bool
}

func (bd DefaultBoardDrawer) DrawBoard(board internal.Board) ShowedGame {
	return ShowedGame{
		Lines: internal.DrawBoardIntoStringArray(board, bd.RevealEverything),
	}
}

func DrawGame(game game.Game, drawer BoardDrawer) (ShowedGame, error) {
	board, err := game.Board()

	if err != nil {
		return ShowedGame{}, err
	}

	return drawer.DrawBoard(board), nil
}

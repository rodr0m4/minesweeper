package operation

import "minesweeper/internal"

type BoardDrawer struct{}

func (bd BoardDrawer) DrawBoardIntoStringArray(board internal.Board, revealEverything bool) []string {
	return internal.DrawBoardIntoStringArray(board, revealEverything)
}

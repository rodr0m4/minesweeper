package game

import (
	"minesweeper/internal"
)

type Fake struct {
	IsStartedFunc func() bool
	BoardFunc     func() (internal.Board, error)
	SyncFunc      func(internal.Board) error
}

func (f Fake) IsStarted() bool {
	return f.IsStartedFunc()
}

func (f Fake) Board() (internal.Board, error) {
	return f.BoardFunc()
}

func (f Fake) Sync(board internal.Board) error {
	return f.SyncFunc(board)
}

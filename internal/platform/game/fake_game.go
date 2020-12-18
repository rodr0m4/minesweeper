package game

import (
	"minesweeper/internal"
)

type Fake struct {
	IsStartedFunc  func() bool
	BoardFunc      func() (internal.Board, error)
	SyncFunc       func(internal.Board) error
	FinishFunc     func()
	IsFinishedFunc func() bool
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

func (f Fake) Finish() {
	f.FinishFunc()
}

func (f Fake) IsFinished() bool {
	return f.IsFinishedFunc()
}

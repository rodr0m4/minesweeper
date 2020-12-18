package game

import (
	"minesweeper/internal"
	"sync"
)

// An InMemory representation of a Game, used for local debugging
type InMemory struct {
	mutex      sync.Mutex
	board      *internal.Board
	isFinished bool
}

func NewInMemory(board *internal.Board) *InMemory {
	return &InMemory{board: board}
}

var _ Game = (*InMemory)(nil)

func (i *InMemory) IsStarted() bool {
	return i.board != nil
}

func (i *InMemory) Board() (internal.Board, error) {
	if !CanPlay(i) {
		return internal.Board{}, internal.NewInvalidOperation("game not yet started")
	}

	var board internal.Board
	i.execute(func() {
		board = *i.board
	})

	return board, nil
}

func (i *InMemory) Sync(board internal.Board) error {
	i.execute(func() {
		i.board = &board
	})
	return nil
}

func (i *InMemory) Finish() {
	i.execute(func() {
		i.isFinished = true
	})
}

func (i *InMemory) IsFinished() (result bool) {
	i.execute(func() {
		result = i.isFinished
	})
	return
}

func (i *InMemory) execute(f func()) {
	i.mutex.Lock()
	f()
	i.mutex.Unlock()
}

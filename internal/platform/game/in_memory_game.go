package game

import (
	"errors"
	"minesweeper/internal"
	"sync"
)

// An InMemory representation of a Game, used for local debugging
type InMemory struct {
	mutex sync.Mutex
	board *internal.Board
}

func NewInMemory(board *internal.Board) *InMemory {
	return &InMemory{board: board}
}

var _ Game = (*InMemory)(nil)

func (i *InMemory) IsStarted() bool {
	return i.board != nil
}

func (i *InMemory) Board() (internal.Board, error) {
	if !i.IsStarted() {
		return internal.Board{}, errors.New("game not yet started")
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

func (i *InMemory) execute(f func()) {
	i.mutex.Lock()
	f()
	i.mutex.Unlock()
}

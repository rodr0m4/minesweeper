package game

import (
	"minesweeper/internal"
	"sync"
)

var (
	ErrGameDoesNotExist = internal.NewInvalidOperation("game does not exist")
)

type ID int

type InMemoryHolder struct {
	mutex     sync.Mutex
	currentID int
	Games     map[ID]Game
}

type Holder interface {
	Insert(Game) (ID, error)
	Get(ID) (Game, error)
	Delete(ID) error
}

func NewInMemoryHolder() *InMemoryHolder {
	return &InMemoryHolder{
		Games: make(map[ID]Game),
	}
}

func (h *InMemoryHolder) Insert(game Game) (ID, error) {
	var id ID

	h.execute(func() {
		id = ID(h.currentID)
		h.currentID++
		h.Games[id] = game
	})

	return id, nil
}

func (h *InMemoryHolder) Get(id ID) (Game, error) {
	var game Game
	var ok bool

	h.execute(func() {
		game, ok = h.Games[id]
	})

	if !ok {
		return nil, ErrGameDoesNotExist
	}

	return game, nil
}

func (h *InMemoryHolder) Delete(id ID) error {
	h.execute(func() {
		delete(h.Games, id)
	})
	return nil
}

func (h *InMemoryHolder) execute(f func()) {
	h.mutex.Lock()
	f()
	h.mutex.Unlock()
}

package game

import (
	"minesweeper/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InMemory_Board(t *testing.T) {
	expect := internal.NewBoard(10, 10, 1)
	game := NewInMemory(&expect)

	actual, err := game.Board()

	assert.NoError(t, err)
	assert.Equal(t, expect, actual)
}

func Test_InMemory_Not_Started_Board(t *testing.T) {
	game := NewInMemory(nil)

	_, err := game.Board()

	assert.Error(t, err, "game not yet started")
}

func Test_InMemory_Sync(t *testing.T) {
	first := internal.NewBoard(10, 10, 1)
	game := NewInMemory(&first)

	second := internal.NewBoard(4, 5, 2)

	assert.NoError(t, game.Sync(second))

	actual, err := game.Board()
	assert.NoError(t, err)
	assert.Equal(t, second, actual)
}

package game

import (
	"minesweeper/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InMemory_Board(t *testing.T) {
	expect := internal.NewBoardFromInitializedMatrix(internal.Matrix{})
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
	first := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile()},
	})
	game := NewInMemory(&first)

	second := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(internal.WithBomb())},
	})

	assert.NoError(t, game.Sync(second))

	actual, err := game.Board()
	assert.NoError(t, err)
	assert.Equal(t, second, actual)
}

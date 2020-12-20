package internal

import (
	"minesweeper/internal/platform/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_A_Board_Is_Created_With_The_Given_Rows_And_Columns(t *testing.T) {
	board := NewBoard(random.Fixed{}, 10, 10, 0)

	assert.Equal(t, 10, board.Rows())
	assert.Equal(t, 10, board.Columns())

	assert.NotEqual(t, 5, board.Rows())
	assert.NotEqual(t, 5, board.Columns())
}

func Test_Board_Position_For_Valid_Positions(t *testing.T) {
	var position Position
	var err error

	board := NewBoard(random.Fixed{}, 10, 10, 0)

	position, err = board.Position(0, 0)
	assert.NoError(t, err)
	assert.Equal(t, Position{Row: 0, Column: 0}, position)

	position, err = board.Position(2, 4)
	assert.NoError(t, err)
	assert.Equal(t, Position{Row: 2, Column: 4}, position)

	position, err = board.Position(9, 9)
	assert.NoError(t, err)
	assert.Equal(t, Position{Row: 9, Column: 9}, position)
}

func Test_Board_Position_For_Invalid_Positions(t *testing.T) {
	var err error

	board := NewBoard(random.Fixed{}, 10, 10, 0)

	_, err = board.Position(-1, -1)
	assert.Error(t, err)

	_, err = board.Position(0, 10)
	assert.Error(t, err)

	_, err = board.Position(5, 324)
	assert.Error(t, err)
}

func Test_Find_On_Empty_Board(t *testing.T) {
	board := NewBoard(random.Fixed{}, 10, 10, 0)

	tile := board.Find(Position{}) // We find the tile on the top-left

	assert.False(t, tile.HasBomb(), "this board does not have bombs")
	assert.Equal(t, HiddenTile{}, tile.State(), "every tile in a new board should be hidden")
}

func Test_Find_On_A_Board_With_Bombs(t *testing.T) {
	rand := random.NewSequence([]int{0, 0, 1, 1})
	board := NewBoard(rand, 2, 2, 2)

	assert.True(t, board.Find(Position{Row: 0, Column: 0}).hasBomb)
	assert.True(t, board.Find(Position{Row: 1, Column: 1}).hasBomb)

	assert.False(t, board.Find(Position{Row: 1, Column: 0}).hasBomb)
	assert.False(t, board.Find(Position{Row: 0, Column: 1}).hasBomb)
}

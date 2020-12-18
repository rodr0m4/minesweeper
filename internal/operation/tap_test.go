package operation

import (
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tap_Should_Fail_If_Game_Board_Fails(t *testing.T) {
	expected := errors.New("oh no")

	_, actual := Tap{}.Tap(gameWhoseBoardFailsWith(expected), 1, 1)

	assert.Equal(t, expected, actual)
}

func Test_Tap_Should_Fail_If_Position_Is_Invalid(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile(internal.WithBomb())},
		{internal.NewTile(), internal.NewTile()},
	})

	_, err := Tap{}.Tap(gameWhoseBoardSucceedsWith(board), 4, 4)

	assert.Error(t, err)
}

func Test_Tap_Should_Tap_On_Tile_If_Valid_Position(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile()},
		{internal.NewTile(), internal.NewTile()},
	})
	tile := board.Find(internal.Position{})

	result, err := Tap{}.Tap(gameWhoseBoardSucceedsWith(board), 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, internal.NothingResult, result)
	assert.Equal(t, internal.RevealedTile{}, tile.State())
}

func Test_Tap_Should_Call_Finisher_If_Result_Is_Explosion(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(internal.WithBomb())},
	})

	var passes bool
	finisher := gameFinisherMock(func(g game.Game) error {
		passes = true
		return nil
	})

	tap := Tap{
		GameFinisher: finisher,
	}

	result, err := tap.Tap(gameWhoseBoardSucceedsWith(board), 0, 0)

	assert.Equal(t, internal.ExplosionResult, result)
	assert.NoError(t, err)
	assert.True(t, passes)
}

// Mocks

func gameWhoseBoardFailsWith(err error) game.Fake {
	return game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return internal.Board{}, err
		},
	}
}

func gameWhoseBoardSucceedsWith(board internal.Board) game.Fake {
	return game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return board, nil
		},
	}
}

type gameFinisherMock func(game.Game) error

func (m gameFinisherMock) Finish(g game.Game) error {
	return m(g)
}

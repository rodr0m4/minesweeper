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

func Test_Tap_Should_Call_Revealer_If_Result_Of_Tapping_Is_Not_An_Explosion(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile()},
		{internal.NewTile(), internal.NewTile()},
	})
	tile := board.Find(internal.Position{})

	var wasCalled bool
	tap := Tap{
		TileRevealer: tileRevealerMock(func(g game.Game, board *internal.Board, position internal.Position) error {
			wasCalled = true
			return nil
		}),
	}

	result, err := tap.Tap(gameWhoseBoardSucceedsWith(board), 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, internal.NothingResult, result)
	assert.Equal(t, internal.RevealedTile{}, tile.State())
	assert.True(t, wasCalled)
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

type tileRevealerMock func(game.Game, *internal.Board, internal.Position) error

func (t tileRevealerMock) RevealAdjacent(game game.Game, board *internal.Board, position internal.Position) error {
	return t(game, board, position)
}

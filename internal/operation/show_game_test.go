package operation

import (
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Should_Fail_With_Error_If_Cant_Find_Board(t *testing.T) {
	g := game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return internal.Board{}, errors.New("oh no")
		},
	}

	sg := ShowGame{
		BoardDrawer: boardDrawerFunc(func(internal.Board, bool) []string {
			t.Fail()
			return nil
		}),
	}

	_, err := sg.ShowGame(g)

	assert.Error(t, err)
}

func Test_Should_Call_BoardDrawer_If_No_Error(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile()},
	})
	result := []string{"Hello", "World"}

	g := game.Fake{
		BoardFunc: func() (internal.Board, error) {
			return board, nil
		},
	}

	sg := ShowGame{
		BoardDrawer: boardDrawerFunc(func(actual internal.Board, _ bool) []string {
			assert.Equal(t, board, actual)
			return result
		}),
	}

	showed, err := sg.ShowGame(g)

	assert.NoError(t, err)
	assert.Equal(t, result, showed.Lines)
}

type boardDrawerFunc func(internal.Board, bool) []string

func (b boardDrawerFunc) DrawBoardIntoStringArray(board internal.Board, revealEverything bool) []string {
	return b(board, revealEverything)
}

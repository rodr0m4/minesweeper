package operation

import (
	"errors"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DeleteGame_Should_Fail_If_Game_Is_Not_Registered(t *testing.T) {
	dg := DeleteGame{
		Holder: game.NewInMemoryHolder(),
	}

	_, err := dg.DeleteGame(game.ID(1))

	assert.Error(t, err)
}

func Test_DeleteGame_Should_Fail_If_Game_Board_Fails(t *testing.T) {
	err := errors.New("oh no")
	g := gameWhoseBoardFailsWith(err)

	holder := game.NewInMemoryHolder()
	id, _ := holder.Insert(g)

	dg := DeleteGame{
		Holder: holder,
	}

	_, actual := dg.DeleteGame(id)

	assert.Equal(t, err, actual)
}

func Test_DeleteGame_Should_Fail_If_Holder_Delete_Fails(t *testing.T) {
	err := errors.New("oh no")
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{})

	holder := gameHolderMock{
		GetFunc: func(id game.ID) (game.Game, error) {
			return gameWhoseBoardSucceedsWith(board), nil
		},
		DeleteFunc: func(id game.ID) error {
			return err
		},
	}

	dg := DeleteGame{
		Holder: holder,
	}

	_, actual := dg.DeleteGame(game.ID(1))

	assert.Equal(t, err, actual)
}

func Test_DeleteGame_Should_Call_DrawBoard_If_Succeeds(t *testing.T) {
	sg := ShowedGame{Lines: []string{"hello", "world"}}

	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{})
	g := gameWhoseBoardSucceedsWith(board)

	holder := game.NewInMemoryHolder()
	id, _ := holder.Insert(g)

	dg := DeleteGame{
		Holder: holder,
		BoardDrawer: boardDrawerMock(func(actual internal.Board) ShowedGame {
			assert.Equal(t, board, actual)
			return sg
		}),
	}

	actual, err := dg.DeleteGame(id)

	assert.NoError(t, err)
	assert.Equal(t, sg, actual)
}

// Mocks

type boardDrawerMock func(internal.Board) ShowedGame

func (m boardDrawerMock) DrawBoard(board internal.Board) ShowedGame {
	return m(board)
}

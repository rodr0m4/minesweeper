package operation

import (
	"errors"
	"fmt"
	"minesweeper/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Mark_Should_Fail_If_Game_Board_Fails(t *testing.T) {
	err := errors.New("oh no")
	game := gameWhoseBoardFailsWith(err)

	assert.Equal(t, err, Mark{}.Mark(game, 0, 0, 0))
}

func Test_Mark_Should_Fail_If_Invalid_Position(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{})
	game := gameWhoseBoardSucceedsWith(board)

	assert.Error(t, Mark{}.Mark(game, 1, 1, 0))
}

func Test_Mark_Should_Fail_If_Marking_Revealed_Tile(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(internal.Reveal())},
	})
	game := gameWhoseBoardSucceedsWith(board)

	assert.Error(t, Mark{}.Mark(game, 0, 0, internal.FlagMark))
}

func Test_Mark_Should_Succeed_And_Sync(t *testing.T) {
	type Case struct {
		name  string
		state internal.TileState
	}

	cases := []Case{
		{
			name:  "hidden",
			state: internal.HiddenTile{},
		},
		{
			name:  "marked",
			state: internal.MarkedTile{},
		},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("mark should succeed and sync on %s tile", tt.name)
		t.Run(name, func(t *testing.T) {
			board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
				{internal.NewTile(internal.WithState(tt.state))},
			})
			game := gameWhoseBoardSucceedsWith(board)
			game.SyncFunc = func(actual internal.Board) error {
				assert.Equal(t, board, actual)
				assert.Equal(t, internal.MarkedTile{Mark: internal.FlagMark}, board.Find(internal.Position{}).State())
				return nil
			}

			assert.NoError(t, Mark{}.Mark(game, 0, 0, internal.FlagMark))
		})
	}
}

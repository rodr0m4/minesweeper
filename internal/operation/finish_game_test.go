package operation

import (
	"errors"
	"minesweeper/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Finish_Game_Fails_When_Game_Board_Fails(t *testing.T) {
	expected := errors.New("oh no")

	actual := FinishGame{}.Finish(gameWhoseBoardFailsWith(expected))

	assert.Equal(t, expected, actual)
}

func Test_Finish_Game_Reveals_Every_Tile_In_The_Board_And_Syncs(t *testing.T) {
	board := internal.NewBoard(1, 1, 0)

	g := gameWhoseBoardSucceedsWith(board)

	var passes bool
	g.SyncFunc = func(board internal.Board) error {
		board.Traverse(func(tile *internal.Tile, position internal.Position) {
			passes = assert.IsType(t, internal.RevealedTile{}, tile.State())
		})
		return nil
	}

	assert.NoError(t, FinishGame{}.Finish(g))
	assert.True(t, passes)
}

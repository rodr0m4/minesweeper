package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StartGame_Already_Started_Should_Fail(t *testing.T) {
	g := game.Fake{
		IsStartedFunc: func() bool {
			return true
		},
	}

	assert.Error(t, StartGame{}.StartGame(g, 0, 0, 0))
}

// TODO: Move this behaviour to Board
func Test_StartGame_With_Invalid_Board(t *testing.T) {
	type Case struct {
		rows    int
		columns int
		bombs   int
	}

	cases := []Case{
		{
			rows:    -1,
			columns: 2,
			bombs:   2,
		},
		{
			rows:    0,
			columns: 2,
			bombs:   2,
		},
		{
			rows:    2,
			columns: -1,
			bombs:   2,
		},
		{
			rows:    2,
			columns: 0,
			bombs:   2,
		},
		{
			rows:    2,
			columns: 2,
			bombs:   -1,
		},
		{
			rows:    2,
			columns: 2,
			bombs:   0,
		},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("rows: %d, columns: %d, bombs: %d", tt.rows, tt.columns, tt.bombs)
		t.Run(name, func(t *testing.T) {
			g := game.Fake{
				IsStartedFunc: func() bool {
					return false
				},
			}

			assert.Error(t, StartGame{}.StartGame(g, tt.rows, tt.columns, tt.bombs))
		})
	}
}

func Test_StartGame_Should_Sync_Valid_Board_With_One_Bomb(t *testing.T) {
	rows := 2
	columns := 2
	bombs := 1

	g := game.Fake{
		IsStartedFunc: func() bool {
			return false
		},
		SyncFunc: func(board internal.Board) error {
			assert.Equal(t, rows, board.Rows())
			assert.Equal(t, columns, board.Columns())
			assert.Equal(t, bombs, board.Bombs())

			return nil
		},
	}

	sg := StartGame{
		Rand: random.NewSequence([]int{0, 0, 1, 1}),
	}

	assert.NoError(t, sg.StartGame(g, rows, columns, bombs))
}

func Test_StartGame_Should_Sync_Valid_Board_With_Game(t *testing.T) {
	rows := 2
	columns := 2
	bombs := 2

	g := game.Fake{
		IsStartedFunc: func() bool {
			return false
		},
		SyncFunc: func(board internal.Board) error {
			assert.Equal(t, rows, board.Rows())
			assert.Equal(t, columns, board.Columns())
			assert.Equal(t, bombs, board.Bombs())

			return nil
		},
	}

	start := StartGame{
		Rand: random.NewSequence([]int{0, 0, 1, 1}),
	}

	assert.NoError(t, start.StartGame(g, rows, columns, bombs))
}

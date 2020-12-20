package operation

import (
	"errors"
	"fmt"
	"minesweeper/internal/platform/game"
	"minesweeper/internal/platform/random"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateGame_Invalid_Board(t *testing.T) {
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
			var cg CreateGame

			_, _, err := cg.CreateGame(tt.rows, tt.columns, tt.bombs)

			assert.Equal(t, invalidBoardError(tt.rows, tt.columns, tt.bombs), err)
		})
	}
}

func Test_CreateGame_Should_Fail_If_Holder_Fails(t *testing.T) {
	err := errors.New("oh no")

	cg := CreateGame{
		Holder: gameHolderMock{
			InsertFunc: func(g game.Game) (game.ID, error) {
				return 2, err
			},
		},
		Rand: random.NewSequence([]int{0, 0, 1, 1}),
	}

	_, _, actual := cg.CreateGame(2, 2, 2)

	assert.Equal(t, err, actual)
}

func Test_CreateGame_Should_Call_Holder(t *testing.T) {
	cg := CreateGame{
		Holder: gameHolderMock{
			InsertFunc: func(g game.Game) (game.ID, error) {
				return 2, nil
			},
		},
		Rand: random.NewSequence([]int{0, 0, 1, 1}),
	}

	id, _, err := cg.CreateGame(2, 2, 2)

	assert.Equal(t, game.ID(2), id)
	assert.NoError(t, err)
}

// Mocks

type gameHolderMock struct {
	InsertFunc func(game.Game) (game.ID, error)
	GetFunc    func(game.ID) (game.Game, error)
}

func (g gameHolderMock) Insert(game game.Game) (game.ID, error) {
	return g.InsertFunc(game)
}

func (g gameHolderMock) Get(id game.ID) (game.Game, error) {
	return g.GetFunc(id)
}

package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RevealAdjacent_Should_Sync_At_The_End(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{}) // Empty board so it is a noop
	g := game.Fake{
		SyncFunc: func(actual internal.Board) error {
			assert.Equal(t, board, actual)
			return nil
		},
	}

	ra := RevealAdjacent{}

	assert.NoError(t, ra.RevealAdjacent(g, &board, internal.Position{}))
}

func Test_RevealAdjacent_In_Corner_Single_Pass(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile()},
		{internal.NewTile(), internal.NewTile(internal.WithBomb())},
	})

	ra := RevealAdjacent{}

	assert.NoError(t, ra.RevealAdjacent(
		newGameThatSucceedsSyncing(),
		&board,
		internal.Position{},
	))

	assertIsRevealed(t, board, internal.Position{Row: 0, Column: 1})
	assertIsRevealed(t, board, internal.Position{Row: 1, Column: 0})
	assertIsNotRevealed(t, board, internal.Position{Row: 1, Column: 1})
}

func Test_RevealAdjacent_In_Corner_Multiple_Passes(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile(internal.WithBomb()), internal.NewTile(internal.WithBomb())},
		{internal.NewTile(), internal.NewTile(), internal.NewTile()},
		{internal.NewTile(internal.WithBomb()), internal.NewTile(internal.WithBomb()), internal.NewTile()},
	})

	ra := RevealAdjacent{}

	assert.NoError(t, ra.RevealAdjacent(
		newGameThatSucceedsSyncing(),
		&board,
		internal.Position{Row: 2, Column: 2},
	))

	assertIsNotRevealed(t, board, internal.Position{Row: 2, Column: 1})
	assertIsNotRevealed(t, board, internal.Position{Row: 2, Column: 0})

	assertIsRevealed(t, board, internal.Position{Row: 1, Column: 2})
	assertIsRevealed(t, board, internal.Position{Row: 1, Column: 1})
	assertIsRevealed(t, board, internal.Position{Row: 1, Column: 0})

	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 2})
	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 1})
	assertIsRevealed(t, board, internal.Position{Row: 0, Column: 0})
}

func Test_RevealAdjacent_In_The_Middle_Complex_Passes(t *testing.T) {
	board := internal.NewBoardFromInitializedMatrix(internal.Matrix{
		{internal.NewTile(), internal.NewTile(), internal.NewTile(internal.WithBomb()), internal.NewTile()},
		{internal.NewTile(internal.WithBomb()), internal.NewTile(internal.WithBomb()), internal.NewTile(), internal.NewTile(internal.Flag())},
		{internal.NewTile(internal.WithBomb()), internal.NewTile(), internal.NewTile(), internal.NewTile(internal.WithBomb())},
		{internal.NewTile(internal.WithBomb()), internal.NewTile(), internal.NewTile(internal.Flag()), internal.NewTile()},
	})

	ra := RevealAdjacent{}

	assert.NoError(t, ra.RevealAdjacent(
		newGameThatSucceedsSyncing(),
		&board,
		internal.Position{Row: 1, Column: 2},
	))

	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 0})
	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 1})
	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 2})
	assertIsNotRevealed(t, board, internal.Position{Row: 0, Column: 3})

	assertIsNotRevealed(t, board, internal.Position{Row: 1, Column: 0})
	assertIsNotRevealed(t, board, internal.Position{Row: 1, Column: 1})
	// 1, 2 is the originator, we do not care about it
	assertIsNotRevealed(t, board, internal.Position{Row: 1, Column: 3})

	assertIsNotRevealed(t, board, internal.Position{Row: 2, Column: 0})
	assertIsRevealed(t, board, internal.Position{Row: 2, Column: 1})
	assertIsRevealed(t, board, internal.Position{Row: 2, Column: 2})
	assertIsNotRevealed(t, board, internal.Position{Row: 2, Column: 3})

	assertIsNotRevealed(t, board, internal.Position{Row: 3, Column: 0})
	assertIsRevealed(t, board, internal.Position{Row: 3, Column: 1})
	assertIsNotRevealed(t, board, internal.Position{Row: 3, Column: 2})
	assertIsNotRevealed(t, board, internal.Position{Row: 3, Column: 3})
}

// Helpers

func assertIsRevealed(t *testing.T, board internal.Board, position internal.Position) {
	assert.IsType(t, internal.RevealedTile{}, board.Find(position).State())
}

func assertIsNotRevealed(t *testing.T, board internal.Board, position internal.Position) {
	state := board.Find(position).State()

	switch state.(type) {
	case internal.RevealedTile:
		assert.Fail(t, fmt.Sprintf("tile at position (%d;%d) was a %T", position.Row, position.Column, state))
	}
}

func newGameThatSucceedsSyncing() game.Fake {
	return game.Fake{
		SyncFunc: func(board internal.Board) error {
			return nil
		},
	}
}

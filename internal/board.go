package internal

import (
	"fmt"
	"minesweeper/internal/platform/random"
)

type Matrix = [][]*Tile

type Board struct {
	matrix  Matrix
	rows    int
	columns int
	bombs   int
}

type Position struct {
	Row    int
	Column int
}

// NewBoardFromInitializedMatrix expects a fully initialized matrix instead of
// generating it using randomness. It is useful for testing and serializing or
// deserializing.
func NewBoardFromInitializedMatrix(matrix Matrix) Board {
	var rows, columns, bombs int

	rows = len(matrix)
	if rows != 0 {
		columns = len(matrix[0])
	}

	traverse(rows, columns, func(position Position) {
		if matrix[position.Row][position.Column].hasBomb {
			bombs++
		}
	})

	return Board{
		matrix:  matrix,
		rows:    rows,
		columns: columns,
		bombs:   bombs,
	}
}

// NewBoard creates a board, first creating the matrix with bombs at random
// places. It is used when first creating the Board.
func NewBoard(rand random.Intn, rows, columns, bombs int) Board {
	matrix := newMatrix(rows, columns)

	for _, position := range bombPositions(rand, rows, columns, bombs) {
		matrix[position.Row][position.Column].hasBomb = true
	}

	return NewBoardFromInitializedMatrix(matrix)
}

func (b Board) Rows() int {
	return b.rows
}

func (b Board) Columns() int {
	return b.columns
}

func (b Board) Bombs() int {
	return b.bombs
}

func (b Board) Find(position Position) *Tile {
	row := position.Row
	column := position.Column

	// Manual bounds check so the error message is better than accessing the
	// matrix with bad indexes
	if b.isOutOfBounds(row, column) {
		panic(b.outOfBoundsError(row, column))
	}

	return b.matrix[row][column]
}

// Traverse goes through the matrix, top-left to bottom right, first every column
// then every row.
func (b Board) Traverse(f func(*Tile, Position)) {
	traverse(b.rows, b.columns, func(position Position) {
		tile := b.Find(position)
		f(tile, position)
	})
}

// Position generates either a valid position for that pair, or an error.
func (b Board) Position(row, column int) (Position, error) {
	if b.isOutOfBounds(row, column) {
		return Position{}, b.outOfBoundsError(row, column)
	}

	return Position{Row: row, Column: column}, nil
}

// FindAdjacentTo gets every adjacent Position to the given one that are
// within bounds of the Board.
func (b Board) FindAdjacentTo(position Position) []Position {
	row := position.Row
	column := position.Column

	potentialPositions := []Position{
		{Row: row - 1, Column: column}, // Top tile
		{Row: row + 1, Column: column}, // Bottom tile
		{Row: row, Column: column - 1}, // Left tile
		{Row: row, Column: column + 1}, // Right tile
	}

	// We filter every out of bounds value
	var adjacent []Position

	for _, position := range potentialPositions {
		if !b.isOutOfBounds(position.Row, position.Column) {
			adjacent = append(adjacent, position)
		}
	}

	return adjacent
}

// Private members

func newMatrix(rows int, columns int) Matrix {
	matrix := make(Matrix, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]*Tile, columns)

		for j := 0; j < columns; j++ {
			matrix[i][j] = NewTile()
		}
	}
	return matrix
}

// TODO: Can this be optimized? maybe generating positions and checking them
// at the same time in two goroutines?
func bombPositions(rand random.Intn, rows, columns, bombs int) []Position {
	var positions []Position

	for i := 0; i < bombs; i++ {
		newPosition := Position{
			Row:    rand.Intn(rows),
			Column: rand.Intn(columns),
		}

		if containsPosition(positions, newPosition) {
			// We still have to keep this iteration
			i--
			continue
		}

		positions = append(positions, newPosition)
	}

	return positions
}

func traverse(rows, columns int, f func(Position)) {
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			f(Position{Row: i, Column: j})
		}
	}
}

func containsPosition(positions []Position, position Position) bool {
	for _, p := range positions {
		if p == position {
			return true
		}
	}
	return false
}

func (b Board) isOutOfBounds(row int, column int) bool {
	return row < 0 || column < 0 || b.rows <= row || b.columns <= column
}

func (b Board) outOfBoundsError(row int, column int) error {
	return NewInvalidOperation(fmt.Sprintf("invalid position access (%d, %d) when Board is %dx%d, use Board::Position for creating a valid Position", row, column, b.rows, b.columns))
}

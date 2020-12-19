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

// This constructor does not use randomness, for testing purposes
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

func NewBoard(rand random.Intn, rows, columns, bombs int) Board {
	matrix := newMatrix(rows, columns)

	for _, position := range bombPositions(rand, rows, columns, bombs) {
		matrix[position.Row][position.Column].hasBomb = true
	}

	return NewBoardFromInitializedMatrix(matrix)
}

func (b Board) HasRows(i int) bool {
	return b.rows == i
}

func (b Board) HasColumns(i int) bool {
	return b.columns == i
}

func (b Board) HasBombs(i int) bool {
	return b.bombs == i
}

func (b Board) Area() int {
	return b.rows * b.columns
}

func (b Board) HasArea(i int) bool {
	return b.Area() == i
}

func (b Board) Find(position Position) *Tile {
	row := position.Row
	column := position.Column

	return b.matrix[row][column]
}

func (b Board) Traverse(f func(*Tile, Position)) {
	traverse(b.rows, b.columns, func(position Position) {
		tile := b.Find(position)
		f(tile, position)
	})
}

func (b Board) Position(row, column int) (Position, error) {
	if b.isOutOfBounds(row, column) {
		return Position{}, b.outOfBoundsError(row, column)
	}

	return Position{Row: row, Column: column}, nil
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

// TODO: Optimize
// Naive algorithm
// TODO: There is a bug with one bomb
func bombPositions(rand random.Intn, rows, columns, bombs int) []Position {
	var positions []Position

	for i := 0; i < bombs-1; i++ {
		found := true
		for found {
			newPosition := Position{
				Row:    rand.Intn(rows),
				Column: rand.Intn(columns),
			}

			for _, position := range positions {
				found = position == newPosition
				if found {
					break
				}
			}

			if len(positions) == 0 || !found {
				positions = append(positions, newPosition)
			}
		}
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

func (b Board) isOutOfBounds(row int, column int) bool {
	return row < 0 || column < 0 || b.rows <= row || b.columns <= column
}

func (b Board) outOfBoundsError(row int, column int) error {
	return NewInvalidOperation(fmt.Sprintf("invalid position access (%d, %d) when Board is %dx%d, use Board::Position for creating a valid Position", row, column, b.rows, b.columns))
}

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

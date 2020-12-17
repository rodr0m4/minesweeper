package internal

import (
	"fmt"
	"math/rand"
)

type Board struct {
	matrix  [][]*Tile
	rows    int
	columns int
	bombs   int
}

type Position struct {
	Row    int
	Column int
}

func NewBoard(rows, columns, bombs int) Board {
	matrix := newMatrix(rows, columns)

	for _, position := range bombPositions(rows, columns, bombs) {
		matrix[position.Row][position.Column].hasBomb = true
	}

	return Board{
		matrix:  matrix,
		rows:    rows,
		columns: columns,
		bombs:   bombs,
	}
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

func newMatrix(rows int, columns int) [][]*Tile {
	matrix := make([][]*Tile, rows)
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
func bombPositions(rows, columns, bombs int) []Position {
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

package internal

import (
	"fmt"
	"strings"
)

func DrawBoardIntoStringArray(board Board, revealEverything bool) []string {
	var rows []string
	sb := new(strings.Builder)

	row := 0
	board.Traverse(func(tile *Tile, position Position) {
		if row != position.Row {
			row++
			flush(sb, &rows)
		}

		sb.WriteRune('|')
		writeTile(sb, tile, revealEverything)
	})

	flush(sb, &rows)

	return rows
}

func writeTile(sb *strings.Builder, tile *Tile, revealEverything bool) {
	if revealEverything {
		writeTileRevealing(sb, tile)
	} else {
		writeTileNotRevealing(sb, tile)
	}
}

func writeTileNotRevealing(sb *strings.Builder, tile *Tile) {
	switch state := tile.State().(type) {
	case MarkedTile:
		var r rune
		if state.Mark == QuestionMark {
			r = '?'
		}
		if state.Mark == FlagMark {
			r = 'F'
		}
		sb.WriteRune(r)
	case HiddenTile:
		sb.WriteRune(' ')
	default:
		if tile.HasBomb() {
			sb.WriteRune('X')
		} else {
			sb.WriteRune('/')
		}
	}
}

func writeTileRevealing(sb *strings.Builder, tile *Tile) {
	sb.WriteRune(bombToRune(tile))
	sb.WriteRune(stateToRune(tile))
}

func stateToRune(tile *Tile) rune {
	switch tile.State().(type) {
	case HiddenTile:
		return 'H'
	case RevealedTile:
		return 'R'
	case MarkedTile:
		return 'F'
	}

	panic(fmt.Errorf("unrecheable code: invalid tile state %d", tile.State()))
}

func bombToRune(tile *Tile) rune {
	if tile.HasBomb() {
		return 'X'
	} else {
		return 'O'
	}
}

func flush(sb *strings.Builder, rows *[]string) {
	sb.WriteRune('|')
	*rows = append(*rows, sb.String())
	sb.Reset()
}

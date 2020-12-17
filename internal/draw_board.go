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
		sb.WriteRune(bombToRun(tile))
		sb.WriteRune(stateToRune(tile))
	})

	flush(sb, &rows)

	return rows
}

func stateToRune(tile *Tile) rune {
	switch tile.State().(type) {
	case HiddenTile:
		return 'H'
	case RevealedTile:
		return 'R'
	case FlaggedTile:
		return 'F'
	}

	panic(fmt.Errorf("unrecheable code: invalid tile state %d", tile.State()))
}

func bombToRun(tile *Tile) rune {
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

package operation

import (
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type RevealAdjacent struct{}

func (ra RevealAdjacent) RevealAdjacent(game game.Game, board *internal.Board, position internal.Position) error {
	ra.revealAdjacent(game, board, position, nil)
	return game.Sync(*board)
}

func (ra RevealAdjacent) revealAdjacent(game game.Game, board *internal.Board, position internal.Position, alreadyRevealed []internal.Position) {
	adjacent := board.FindAdjacentTo(position)

	for _, position := range adjacent {
		if isAlreadyRevealed(position, alreadyRevealed) {
			continue
		}

		tile := board.Find(position)

		if tile.HasBomb() {
			alreadyRevealed = append(alreadyRevealed, position)
			continue
		}

		switch tile.State().(type) {
		case internal.HiddenTile:
			tile.Modify(internal.Reveal())
			alreadyRevealed = append(alreadyRevealed, position)
			ra.revealAdjacent(game, board, position, alreadyRevealed)
		}
	}
}

func isAlreadyRevealed(position internal.Position, alreadyRevealed []internal.Position) bool {
	for _, already := range alreadyRevealed {
		if already == position {
			return true
		}
	}

	return false
}

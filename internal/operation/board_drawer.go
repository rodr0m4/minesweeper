package operation

import (
	"fmt"
	"minesweeper/internal"
	"minesweeper/internal/platform/game"
)

type ShowedGame struct {
	Tiles []ShowedTile `json:"tiles"`
	Lines []string     `json:"lines"`
}

type ShowedTile struct {
	Row     int         `json:"row"`
	Column  int         `json:"column"`
	HasBomb bool        `json:"has_bomb,omitempty"`
	State   ShowedState `json:"state"`
}

type ShowedState string

type BoardDrawer interface {
	DrawBoard(internal.Board) ShowedGame
}

type DefaultBoardDrawer struct {
	ShowLines        bool
	ShowTiles        bool
	RevealEverything bool
}

func (bd DefaultBoardDrawer) DrawBoard(board internal.Board) ShowedGame {
	return ShowedGame{
		Tiles: bd.drawTiles(board),
		Lines: bd.drawLines(board),
	}
}

func DrawGame(game game.Game, drawer BoardDrawer) (ShowedGame, error) {
	board, err := game.Board()

	if err != nil {
		return ShowedGame{}, err
	}

	return drawer.DrawBoard(board), nil
}

func (bd DefaultBoardDrawer) drawLines(board internal.Board) []string {
	if !bd.ShowLines {
		return nil
	}

	return internal.DrawBoardIntoStringArray(board, bd.RevealEverything)
}

func (bd DefaultBoardDrawer) drawTiles(board internal.Board) []ShowedTile {
	if !bd.ShowTiles {
		return nil
	}

	var tiles []ShowedTile

	board.Traverse(func(tile *internal.Tile, position internal.Position) {
		tiles = append(tiles, bd.drawTile(tile, position))
	})

	return tiles
}

func (bd DefaultBoardDrawer) drawTile(tile *internal.Tile, position internal.Position) ShowedTile {
	return ShowedTile{
		Row:     position.Row,
		Column:  position.Column,
		HasBomb: tile.HasBomb(),
		State:   tileStateToShowedState(tile.State()),
	}
}

func tileStateToShowedState(state internal.TileState) ShowedState {
	switch s := state.(type) {
	case internal.HiddenTile:
		return "HIDDEN"
	case internal.MarkedTile:
		switch s.Mark {
		case internal.FlagMark:
			return "FLAG"
		case internal.QuestionMark:
			return "QUESTION_MARK"
		}
	case internal.RevealedTile:
		return "REVEALED"
	}

	panic(fmt.Sprintf("invalid state %T", state))
}

func (bd DefaultBoardDrawer) marshallBomb(tile *internal.Tile) bool {
	hasBomb := tile.HasBomb()

	switch tile.State().(type) {
	case internal.RevealedTile:
		return hasBomb
	}

	return bd.RevealEverything && hasBomb
}

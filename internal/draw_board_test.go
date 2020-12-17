package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Draw_Board_Revealing(t *testing.T) {
	type Case struct {
		name   string
		given  [][]*Tile
		expect []string
	}

	cases := []Case{
		{
			name: "on a small board with one bomb",
			given: [][]*Tile{
				{&Tile{hasBomb: true}, &Tile{}},
				{&Tile{}, &Tile{}},
			},
			expect: []string{"|XH|OH|", "|OH|OH|"},
		},
		{
			name: "on a bigger board with multiple bombs",
			given: [][]*Tile{
				{&Tile{}, &Tile{}, &Tile{hasBomb: true}, &Tile{}},
				{&Tile{hasBomb: true}, &Tile{}, &Tile{hasBomb: true}, &Tile{hasBomb: true}},
				{&Tile{hasBomb: true}, &Tile{}, &Tile{}, &Tile{hasBomb: true}},
			},
			expect: []string{"|OH|OH|XH|OH|", "|XH|OH|XH|XH|", "|XH|OH|OH|XH|"},
		},
		{
			name: "with flagged and revealed tiles",
			given: [][]*Tile{
				{&Tile{hasBomb: true, state: FlaggedTile}, &Tile{}},
				{&Tile{hasBomb: true}, &Tile{state: RevealedTile}},
			},
			expect: []string{"|XF|OH|", "|XH|OR|"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			board := Board{
				matrix:  tt.given,
				rows:    len(tt.given),
				columns: len(tt.given[0]), // TODO: this is unsafe
			}

			assert.Equal(t, tt.expect, DrawBoardIntoStringArray(board, true))
		})
	}
}

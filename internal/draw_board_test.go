package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Draw_Board_Revealing(t *testing.T) {
	type Case struct {
		name   string
		given  Matrix
		expect []string
	}

	cases := []Case{
		{
			name: "on a small board with one bomb",
			given: Matrix{
				{NewTile(WithBomb()), NewTile()},
				{NewTile(), NewTile()},
			},
			expect: []string{"|XH|OH|", "|OH|OH|"},
		},
		{
			name: "on a bigger board with multiple bombs",
			given: Matrix{
				{NewTile(), NewTile(), NewTile(WithBomb()), NewTile()},
				{NewTile(WithBomb()), NewTile(), NewTile(WithBomb()), NewTile(WithBomb())},
				{NewTile(WithBomb()), NewTile(), NewTile(), NewTile(WithBomb())},
			},
			expect: []string{"|OH|OH|XH|OH|", "|XH|OH|XH|XH|", "|XH|OH|OH|XH|"},
		},
		{
			name: "with flagged and revealed tiles",
			given: Matrix{
				{NewTile(WithBomb(), WithMark(FlagMark)), NewTile()},
				{NewTile(WithBomb()), NewTile(Reveal())},
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

func Test_Draw_Board_Not_Revealing(t *testing.T) {
	board := NewBoardFromInitializedMatrix(Matrix{
		{NewTile(WithBomb()), NewTile()},
		{NewTile(WithBomb(), Reveal()), NewTile(Reveal())},
		{NewTile(WithBomb(), WithMark(QuestionMark)), NewTile(WithBomb(), WithMark(FlagMark))},
		{NewTile(WithMark(QuestionMark)), NewTile(WithMark(FlagMark))},
	})

	expect := []string{
		"| | |",
		"|X|/|",
		"|?|F|",
		"|?|F|",
	}

	assert.Equal(t, expect, DrawBoardIntoStringArray(board, false))
}

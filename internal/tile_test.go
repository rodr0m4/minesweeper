package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Tile_Is_Created_Without_Bomb_And_Hidden_By_Default(t *testing.T) {
	tile := NewTile()

	assert.False(t, tile.HasBomb())
	assert.Equal(t, HiddenTile{}, tile.State())
}

func Test_Tap_On_Hidden_Tile_Without_Bomb_Reveals_It(t *testing.T) {
	var tile Tile

	result, err := tile.Tap()

	assert.NoError(t, err)
	assert.Equal(t, NothingResult, result)
	assert.Equal(t, RevealedTile{}, tile.State())
}

func Test_Tap_On_Hidden_Tile_With_Bomb_Explodes_And_Reveals_It(t *testing.T) {
	tile := NewTile(WithBomb())

	result, err := tile.Tap()

	assert.NoError(t, err)
	assert.Equal(t, ExplosionResult, result)
	assert.Equal(t, RevealedTile{}, tile.State())
}

func Test_Tap_Invalid_Operation(t *testing.T) {
	type Case struct {
		name              string
		state             TileState
		showableTileState string
	}

	cases := []Case{
		{
			name:              "in a revealed tile",
			state:             RevealedTile{},
			showableTileState: "revealed",
		},
		{
			name:              "in a flagged tile",
			state:             FlaggedTile{},
			showableTileState: "flagged",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			originalTile := NewTile(WithState(tt.state))
			copiedTile := originalTile // So we avoid comparing field by field

			_, err := originalTile.Tap()

			assert.Error(t, err, fmt.Sprintf("invalid operation: can not tap on %s tiles", tt.showableTileState))
			assert.Equal(t, copiedTile, originalTile)
		})
	}
}

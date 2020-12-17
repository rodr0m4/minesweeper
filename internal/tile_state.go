package internal

// TileState is a sealed union type
type TileState interface {
	// Marker method
	isTileState()
}

// HidenTileState is when a tile is not revealed, it can be tagged with a
// number of adjacent tiles
type HiddenTile struct {
	// How many Tiles are close to this one
	Adjacent int
}

func (h HiddenTile) isTileState() {
}

type FlaggedTile struct{}

func (f FlaggedTile) isTileState() {
}

type RevealedTile struct{}

func (r RevealedTile) isTileState() {
}

package internal

// TileState is a sealed union type
type TileState interface {
	// Marker method
	isTileState()
}

type HiddenTile struct {
}

func (h HiddenTile) isTileState() {
}

type MarkedTile struct {
	Mark TileMark
}

type TileMark int

const (
	FlagMark TileMark = iota
	QuestionMark
)

func (f MarkedTile) isTileState() {
}

type RevealedTile struct{}

func (r RevealedTile) isTileState() {
}

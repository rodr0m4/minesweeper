package internal

var (
	ErrCantTapOnRevealedTile = NewInvalidOperation("can not tap on revealed tile")
	ErrCantTapOnFlaggedTile  = NewInvalidOperation("can not tap on flagged tile")
)

// Tiles have private state because some combination of attributes could make it illegal (revealed + flagged, for
// example).
type Tile struct {
	hasBomb bool
	state   TileState
}

type NewTileOption func(*Tile)

func WithBomb() NewTileOption {
	return func(tile *Tile) {
		tile.hasBomb = true
	}
}

func ThatIsFlagged() NewTileOption {
	return WithState(FlaggedTile{})
}

func ThatIsMarked(adjacent int) NewTileOption {
	return WithState(HiddenTile{Adjacent: adjacent})
}

func ThatIsRevealed() NewTileOption {
	return WithState(RevealedTile{})
}

func WithState(state TileState) NewTileOption {
	return func(tile *Tile) {
		tile.state = state
	}
}

func NewTile(opts ...NewTileOption) *Tile {
	tile := &Tile{
		state: HiddenTile{},
	}

	for _, opt := range opts {
		opt(tile)
	}

	return tile
}

func (t *Tile) HasBomb() bool {
	return t.hasBomb
}

func (t *Tile) State() TileState {
	return t.state
}

func (t *Tile) Tap() (result TapResult, err error) {
	switch t.state.(type) {
	case RevealedTile:
		err = ErrCantTapOnRevealedTile
		return
	case FlaggedTile:
		err = ErrCantTapOnFlaggedTile
		return
	}

	t.state = RevealedTile{}

	if t.hasBomb {
		result = ExplosionResult
	}

	return
}

type TapResult int

const (
	NothingResult TapResult = iota
	ExplosionResult
)

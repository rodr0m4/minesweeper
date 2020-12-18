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

type TileOption func(*Tile)

func WithBomb() TileOption {
	return func(tile *Tile) {
		tile.hasBomb = true
	}
}

func Flag() TileOption {
	return WithState(FlaggedTile{})
}

func Mark(adjacent int) TileOption {
	return WithState(HiddenTile{Adjacent: adjacent})
}

func Reveal() TileOption {
	return WithState(RevealedTile{})
}

func WithState(state TileState) TileOption {
	return func(tile *Tile) {
		tile.state = state
	}
}

func NewTile(opts ...TileOption) *Tile {
	tile := &Tile{
		state: HiddenTile{},
	}

	tile.Modify(opts...)

	return tile
}

func (t *Tile) HasBomb() bool {
	return t.hasBomb
}

func (t *Tile) State() TileState {
	return t.state
}

func (t *Tile) Modify(opts ...TileOption) {
	for _, opt := range opts {
		opt(t)
	}
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

func (tr TapResult) String() string {
	return []string{"Nothing", "Explosion"}[tr]
}

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

func (t *Tile) HasBomb() bool {
	return t.hasBomb
}

func (t *Tile) State() TileState {
	return t.state
}

func (t *Tile) Tap() (result TapResult, err error) {
	if t.state == RevealedTile {
		err = ErrCantTapOnRevealedTile
		return
	}

	if t.state == FlaggedTile {
		err = ErrCantTapOnFlaggedTile
		return
	}

	t.state = RevealedTile

	if t.hasBomb {
		result = ExplosionResult
	}

	return
}

type TileState int

const (
	HiddenTile TileState = iota
	RevealedTile
	FlaggedTile
)

type TapResult int

const (
	NothingResult TapResult = iota
	ExplosionResult
)

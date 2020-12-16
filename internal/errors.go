package internal

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
)

func NewInvalidOperation(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidOperation, message)
}

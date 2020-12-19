package random

import "math/rand"

type Real struct{}

func (rr Real) Intn(n int) int {
	return rand.Intn(n)
}

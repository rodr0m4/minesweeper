package random

type Fixed struct {
	N int
}

func (fr Fixed) Intn(n int) int {
	return fr.N
}

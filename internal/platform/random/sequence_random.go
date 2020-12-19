package random

type Sequence struct {
	current  int
	sequence []int
}

func NewSequence(sequence []int) *Sequence {
	return &Sequence{sequence: sequence}
}

func (s *Sequence) Intn(n int) int {
	i := s.sequence[s.current]
	s.current++
	return i
}

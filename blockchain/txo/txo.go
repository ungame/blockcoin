package txo

type Map struct {
	tx2out map[string][]int
}

func NewMap() *Map {
	return &Map{tx2out: make(map[string][]int)}
}

func (s *Map) Set(txID string, outputIndex int) {
	s.tx2out[txID] = append(s.tx2out[txID], outputIndex)
}

func (s *Map) Get(txID string) []int {
	return s.tx2out[txID]
}

func (s *Map) GetAll() map[string][]int {
	return s.tx2out
}

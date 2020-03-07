package main

type Serial struct {
	value int
}

func (s *Serial) Next() int {
	s.value++
	return s.value
}

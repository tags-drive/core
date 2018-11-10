package aggregation

type stack struct {
	data []bool
	len  int
}

func (s *stack) push(b bool) {
	s.data = append(s.data, b)
	s.len++
}

func (s stack) top() bool {
	return s.data[s.len-1]
}

func (s *stack) pop() bool {
	b := s.data[s.len-1]
	s.data = s.data[:s.len-1]
	s.len--
	return b
}

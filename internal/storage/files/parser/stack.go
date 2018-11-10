package parser

type stack struct {
	data []byte
	len  int
}

func (s *stack) push(b byte) {
	s.data = append(s.data, b)
	s.len++
}

func (s stack) top() byte {
	return s.data[s.len-1]
}

func (s *stack) pop() byte {
	b := s.data[s.len-1]
	s.data = s.data[:s.len-1]
	s.len--
	return b
}

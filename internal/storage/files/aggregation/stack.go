package aggregation

// logicalStack is used for parsing logical expression
type logicalStack struct {
	data []byte
	len  int
}

func (s *logicalStack) push(b byte) {
	s.data = append(s.data, b)
	s.len++
}

func (s logicalStack) top() byte {
	return s.data[s.len-1]
}

func (s *logicalStack) pop() byte {
	b := s.data[s.len-1]
	s.data = s.data[:s.len-1]
	s.len--
	return b
}

// processingStack is used for computing parsed logical expression
type processingStack struct {
	data []bool
	len  int
}

func (s *processingStack) push(b bool) {
	s.data = append(s.data, b)
	s.len++
}

func (s processingStack) top() bool {
	return s.data[s.len-1]
}

func (s *processingStack) pop() bool {
	b := s.data[s.len-1]
	s.data = s.data[:s.len-1]
	s.len--
	return b
}

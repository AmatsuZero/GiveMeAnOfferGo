package design_an_ordered_stream

type OrderedStream struct {
	ptr    int
	stream []string
}

func Constructor(n int) OrderedStream {
	return OrderedStream{
		stream: make([]string, n+1),
		ptr:    1,
	}
}

func (s *OrderedStream) Insert(idKey int, value string) []string {
	s.stream[idKey] = value
	start := s.ptr
	for s.ptr < len(s.stream) && s.stream[s.ptr] != "" {
		s.ptr += 1
	}
	return s.stream[start:s.ptr]
}

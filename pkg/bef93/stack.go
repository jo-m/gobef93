package bef93

import "sync"

// TODO: optimize, keep same slice and separate stack pointer
type stack struct {
	s []int64

	//lint:ignore U1000 ignore unused copy guard
	noCopy sync.Mutex
}

func (s *stack) push(val int64) {
	s.s = append(s.s, val)
}

func (s *stack) pop() int64 {
	l := len(s.s)
	if l < 1 {
		return 0
	}

	ret := s.s[l-1]
	s.s = s.s[:l-1]

	return ret
}

func (s *stack) pop2() (int64, int64) {
	return s.pop(), s.pop()
}

func (s *stack) clone() stack {
	arr := make([]int64, len(s.s))
	copy(arr, s.s)

	return stack{
		s: arr,
	}
}

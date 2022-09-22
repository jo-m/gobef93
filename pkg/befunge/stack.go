package befunge

import (
	"errors"
)

type stack struct {
	s []int
}

func (s *stack) push(val int) {
	s.s = append(s.s, val)
	// log.Println("push, stack is now", s.s)
}

func (s *stack) pop() (int, error) {
	l := len(s.s)
	if l < 1 {
		return 0, errors.New("pop on empty stack")
	}

	ret := s.s[l-1]
	s.s = s.s[:l-1]

	// log.Println("pop, val / stack", ret, s.s)

	return ret, nil
}

func (s *stack) pop2() (int, int, error) {
	l := len(s.s)
	if l < 2 {
		return 0, 0, errors.New("pop2 on empty stack")
	}

	ret0, err := s.pop()
	if err != nil {
		panic("should not happen")
	}

	ret1, err := s.pop()
	if err != nil {
		panic("should not happen")
	}

	return ret0, ret1, nil
}

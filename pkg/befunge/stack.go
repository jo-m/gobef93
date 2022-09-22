package befunge

// TODO: optimize, keep same slice and separate stack pointer
type stack struct {
	s []int
}

func (s *stack) push(val int) {
	s.s = append(s.s, val)
	// log.Println("push, stack is now", s.s)
}

func (s *stack) pop() int {
	l := len(s.s)
	if l < 1 {
		return 0
	}

	ret := s.s[l-1]
	s.s = s.s[:l-1]

	// log.Println("pop, val / stack", ret, s.s)

	return ret
}

func (s *stack) pop2() (int, int) {
	a, b := s.pop(), s.pop()
	return a, b
}

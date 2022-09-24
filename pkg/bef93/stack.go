package bef93

import (
	"sync"
)

const initSize = 1 << 6

type stack struct {
	s  []int64
	sp int

	//lint:ignore U1000 ignore unused copy guard
	noCopy sync.Mutex
}

func (s *stack) push(val int64) {
	if len(s.s) == 0 {
		s.s = make([]int64, initSize)
	}

	if len(s.s) == s.sp {
		s.s = append(s.s, make([]int64, len(s.s))...)
	}

	s.s[s.sp] = val
	s.sp++
}

func (s *stack) pop() int64 {
	s.sp--
	if s.sp < 0 {
		s.sp = 0
		return 0
	}

	return s.s[s.sp]
}

func (s *stack) pop2() (int64, int64) {
	return s.pop(), s.pop()
}

func (s *stack) clone() stack {
	arr := make([]int64, len(s.s))
	copy(arr, s.s)

	return stack{
		s:  arr,
		sp: s.sp,
	}
}

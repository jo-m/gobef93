package bef93

import (
	"testing"
)

func Test_stack_Simple(t *testing.T) {
	s := stack{}

	s.push(123)
	for i := int64(1); i <= 100; i++ {
		s.push(i)
	}

	if s.pop() != 100 {
		t.Fatal("invalid value")
	}

	for i := 0; i <= 98; i++ {
		s.pop()
	}

	if s.pop() != 123 {
		t.Fatal("invalid value")
	}

	if s.pop() != 0 {
		t.Fatal("invalid value")
	}

	s.push(456)

	if s.pop() != 456 {
		t.Fatal("invalid value")
	}
}

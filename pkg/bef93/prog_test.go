package bef93

import (
	"errors"
	"testing"
)

func Test_NewProg_Simple(t *testing.T) {
	const code = `>              v
v  ,,,,,"Hello"<
>48*,          v
v,,,,,,"World!"<
>25*,@
`
	prog, err := NewProg(code, Opts{})

	if prog == nil {
		t.Fatalf("prog is nil")
	}

	if err != nil {
		t.Fatalf("err is not is nil: %s", err)
	}

	if prog.w != Width {
		t.Fatalf("invalid w")
	}

	if prog.h != Height {
		t.Fatalf("invalid h")
	}
}

func Test_NewProg_Unicode(t *testing.T) {
	const code = `>              v
v  ,,,,,"Hello"<
>48*,          v
v,,,,,,"Wörld!"<
>25*,@
`
	_, err := NewProg(code, Opts{})

	if !errors.Is(err, ErrNotASCII) {
		t.Fatalf("should fail")
	}

	prog, err := NewProg(code, Opts{AllowUnicode: true})

	if prog == nil {
		t.Fatalf("prog is nil")
	}

	if err != nil {
		t.Fatalf("err is not is nil: %s", err)
	}
}

func Test_NewProg_Size(t *testing.T) {
	const code = `>              v
v  ,,,,,"Hello"<
>48*,          v
v,,,,,,"World!"<                                                                                                       #
>25*,@
























#
`
	_, err := NewProg(code, Opts{})

	if !errors.Is(err, ErrTooLarge) {
		t.Fatalf("should fail")
	}

	prog, err := NewProg(code, Opts{AllowArbitraryCodeSize: true})

	if prog == nil {
		t.Fatalf("prog is nil")
	}

	if err != nil {
		t.Fatalf("err is not is nil: %s", err)
	}

	if prog.w != 120 {
		t.Fatalf("invalid width %d", prog.w)
	}
	if prog.h != 31 {
		t.Fatalf("invalid height %d", prog.h)
	}
}

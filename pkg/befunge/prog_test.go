package befunge

import (
	"os"
	"testing"
)

var code = `
>              v
v  ,,,,,"Hello"<
>48*,          v
v,,,,,,"World!"<
>25*,@
`

func Test_NewProg(t *testing.T) {

	prog := NewProg(code)

	if prog == nil {
		t.Fatalf("prog is nil")
	}
}

func Test_Exec(t *testing.T) {

	prog := NewProg(code)
	if prog == nil {
		t.Fatalf("prog is nil")
	}

	err := prog.Exec(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		t.Fatalf(err.Error())
	}

}

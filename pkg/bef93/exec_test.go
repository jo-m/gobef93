package bef93

import (
	"bytes"
	"testing"
)

func Test_Exec_HelloWorld(t *testing.T) {
	const code = ` >25*"!dlrow ,olleH":v
                  v:,_@
                  >  ^`

	prog, err := NewProg(code, Opts{})
	if prog == nil {
		t.Fatalf("prog is nil")
	}

	if err != nil {
		t.Fatalf(err.Error())
	}

	stdin, stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	proc := NewProc(prog, stdin, stdout, stderr)

	err = proc.Exec()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if stdout.String() != "Hello, world!\n" {
		t.Fatal("should be equal")
	}
}

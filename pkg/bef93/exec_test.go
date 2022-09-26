package bef93

import (
	"bytes"
	"strings"
	"testing"
)

// usage: proc, stdin, stdout, stderr := createProc(t, code)
func createProc(t *testing.T, code string, opts Opts) (*Proc, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	prog, err := NewProg(code, opts)
	if prog == nil {
		t.Fatalf("prog is nil")
	}
	if err != nil {
		t.Fatalf(err.Error())
	}

	stdin, stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
	proc := NewProc(prog, stdin, stdout, stderr)

	return proc, stdin, stdout, stderr
}

func exec2out(t *testing.T, code string, opts Opts, in string) (string, string, error) {
	proc, stdin, stdout, stderr := createProc(t, code, opts)

	stdin.Write([]byte(in))

	err := proc.Exec()
	return stdout.String(), stderr.String(), err
}

func Test_Exec_HelloWorld(t *testing.T) {
	const code = ` >25*"!dlrow ,olleH":v
                  v:,_@
                  >  ^`

	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "Hello, world!\n" {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Add(t *testing.T) {
	out, _, err := exec2out(t, `12+.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "3 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Sub(t *testing.T) {
	out, _, err := exec2out(t, `32-.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "1 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Mul(t *testing.T) {
	out, _, err := exec2out(t, `32*.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "6 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Div(t *testing.T) {
	out, _, err := exec2out(t, `82/.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "4 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Div0_Ask(t *testing.T) {
	out, outErr, err := exec2out(t, `80/.@`, Opts{}, "12")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "12 " {
		t.Fatal("should be equal")
	}
	if outErr != "What do you want 8/0 to be?\n" {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Div0_Fail(t *testing.T) {
	_, _, err := exec2out(t, `80/.@`, Opts{DisallowDivZero: true}, "12")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_Exec_Div0_IoErr(t *testing.T) {
	out, _, err := exec2out(t, `80/.@`, Opts{}, "")

	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "0 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Div0_IoErrFail(t *testing.T) {
	_, _, err := exec2out(t, `80/.@`, Opts{TerminateOnIOErr: true}, "")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_Exec_Mod(t *testing.T) {
	out, _, err := exec2out(t, `73%.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if out != "1 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Mod0_Fail(t *testing.T) {
	_, _, err := exec2out(t, `80%.@`, Opts{}, "")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_Exec_Not(t *testing.T) {
	out, _, err := exec2out(t, `7!.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "0 " {
		t.Fatal("should be equal")
	}

	out, _, err = exec2out(t, `1!.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "0 " {
		t.Fatal("should be equal")
	}

	out, _, err = exec2out(t, `0!.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "1 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Gt(t *testing.T) {
	out, _, err := exec2out(t, "21`.@", Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "1 " {
		t.Fatal("should be equal")
	}

	out, _, err = exec2out(t, "12`.@", Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "0 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Arrows(t *testing.T) {
	code := strings.TrimSpace(`
>  1  v
@     2
.
.
.
^  3  <
	`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "3 2 1 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Pop(t *testing.T) {
	out, _, err := exec2out(t, "123$$$.@", Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "0 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_StrWtrChr(t *testing.T) {
	out, _, err := exec2out(t, `"olleh",,,,,@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "hello" {
		t.Fatal("should be equal")
	}
}

func Test_Exec_WtrInt(t *testing.T) {
	out, _, err := exec2out(t, `999**.@`, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "729 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Rif(t *testing.T) {
	code := strings.TrimSpace(`
>  1  v
  v0.2_3.@
  _7.@
`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "2 7 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Dif(t *testing.T) {
	code := strings.TrimSpace(`
v
             @
     >   0   |
     .       8
     2       .
>  1 |       @
     3
     .
     @

`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "2 8 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Skip(t *testing.T) {
	code := strings.TrimSpace(`
>   1  #23 .. v
v     .. 89#  <
#   @
4   .
5   7
.   6
.   #
>   ^
`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "3 1 8 0 5 0 7 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Put(t *testing.T) {
	code := strings.TrimSpace(`
432pv
v   <
> "3" ..@
`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "4 0 " {
		t.Fatal("should be equal")
	}
}

func Test_Exec_Get(t *testing.T) {
	code := strings.TrimSpace(`
83g,@


        7
`)
	out, _, err := exec2out(t, code, Opts{}, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if out != "7" {
		t.Fatal("should be equal")
	}
}

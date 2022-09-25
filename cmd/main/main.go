package main

import (
	"fmt"
	"os"

	"github.com/jo-m/gobefunge/pkg/bef93"
)

const code = ":0g,:\"~\"`#@_1+0\"Quines are Fun\">_"

func main() {
	prog, err := bef93.NewProg(code, bef93.Opts{})
	if err != nil {
		panic(err)
	}
	fmt.Println(prog.String())

	proc := bef93.NewProc(prog, os.Stdin, os.Stdout, os.Stderr)

	err = proc.Exec()
	if err != nil {
		panic(err)
	}
}

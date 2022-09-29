/*
Package main is a simple CLI frontend for the gobef93 Befunge interpreter.
*/
package main

import (
	"fmt"
	"os"

	"github.com/jo-m/gobefunge/pkg/bef93"
)

const code = `&&..@`

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

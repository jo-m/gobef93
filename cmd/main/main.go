package main

import (
	"io"
	"log"
	"os"

	"github.com/jo-m/gobefunge/pkg/befunge"
)

const code = `
:0g,:"~"]#@_1+0"Quines are Fun">_
`

func main() {
	prog := befunge.NewProg93(code)
	if prog == nil {
		panic("prog is nil")
	}

	// fmt.Println(prog.String())

	log.SetOutput(io.Discard)

	err := prog.Exec(os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		panic(err)
	}
}

package main

import (
	"io"
	"log"
	"os"

	"github.com/jo-m/gobefunge/pkg/bef93"
)

const code = `
:0g,:"~"]#@_1+0"Quines are Fun">_
`

func main() {
	prog := bef93.NewProg(code, bef93.Opts{})
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

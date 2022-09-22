package main

import (
	"io"
	"log"
	"os"

	"github.com/jo-m/gobefunge/pkg/befunge"
)

const code = `
 >25*"!dlrow ,olleH":v
                  v:,_@
                  >  ^
`

func main() {
	prog := befunge.NewProg(code)
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

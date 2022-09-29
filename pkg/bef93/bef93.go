/*
Package bef93 implements a Befunge93 interpreter.

Sample usage:

	prog, err := bef93.NewProg(code, bef93.Opts{})

	if err != nil {
		panic(err)
	}

	proc := bef93.NewProc(prog, os.Stdin, os.Stdout, os.Stderr)
	err = proc.Exec()

	if err != nil {
		panic(err)
	}
*/
package bef93

/*
Package main is a simple CLI frontend for the gobef93 Befunge interpreter.
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jo-m/gobef93/pkg/bef93"
)

type mainOpts struct {
	printProg bool
}

func mustParseFlags() (string, bef93.Opts, mainOpts) {
	opts := bef93.Opts{}

	flag.BoolVar(&opts.ReadErrorUndefined, "read_error_undefined", false, "If true, & will push an undefined number to stack instead of -1. Befunge 93 standard option.")
	flag.BoolVar(&opts.IgnoreUnsupportedInstructions, "ignore_unsupported_instructions", false, "If true, unsupported instructions will be ignored. Befunge 93 standard option.")

	flag.BoolVar(&opts.AllowArbitraryCodeSize, "allow_arbitrary_code_size", false, "Allow code of arbitrary size, code smaller than standard size will be padded to standard size. Non standard option.")
	flag.BoolVar(&opts.AllowUnicode, "allow_unicode", false, "Allow unicode in the interpreted code. Non standard option.")
	flag.BoolVar(&opts.DisallowDivZero, "disallow_div_zero", false, "Terminate on division by 0. Non standard option.")
	flag.Int64Var(&opts.RandSeed, "rand_seed", 0, "Fixed random seed. If 0, the generator is seeded randomly internally. Non standard option.")
	flag.BoolVar(&opts.TerminateOnIOErr, "terminate_on_io_err", false, "Terminate on I/O errors instead of ignoring them. Non standard option.")
	flag.BoolVar(&opts.TerminateOnPutGetOutOfBounds, "terminate_on_put_get_out_of_bounds", false, "Terminate if a 'g' or 'p' operation is out of bounds, instead of pushing 0 or discading the pop() value. Non standard option.")

	mainOpts := mainOpts{}

	flag.BoolVar(&mainOpts.printProg, "print_prog", false, "Print program grid to stderr before execution. Non standard option.")

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(w, `Executes a Befunge-93 program file.
Takes a single positional argument, which is the file to execute.
For more details on the options, see the docstrings on the bef93.Opts struct.`+"\n")

		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "missing positional argument (file name)\n")
		flag.Usage()
		os.Exit(1)
	}

	return flag.Arg(0), opts, mainOpts
}

func mustGetCode(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	code, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	return string(code)
}

func main() {
	srcFile, opts, mainOpts := mustParseFlags()
	code := mustGetCode(srcFile)

	prog, err := bef93.NewProg(code, opts)
	if err != nil {
		panic(err)
	}

	if mainOpts.printProg {
		fmt.Fprintln(os.Stderr, prog.String())
	}

	proc := bef93.NewProc(prog, os.Stdin, os.Stdout, os.Stderr)
	err = proc.Exec()
	if err != nil {
		panic(err)
	}
}

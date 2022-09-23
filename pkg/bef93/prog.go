package bef93

import (
	"strings"
	"unicode"
)

const (
	Width  = 80
	Height = 25
)

// Opts contains supported options.
// See https://github.com/catseye/Befunge-93/blob/master/src/bef.c#L46.
// Zero value is good to use and represents the default options.
type Opts struct {
	NoFixOffByOne                 bool // TODO: implement
	ReadErrorUndefined            bool // TODO: implement
	IgnoreUnsupportedInstructions bool // TODO: implement
	WrapLongLines                 bool // TODO: implement
	WrapHashInconsistently        bool // TODO: implement

	// non-standard options
	AllowInfiniteCodeSize bool // TODO: implement
	AllowUnicode          bool // TODO: implement
}

type Prog struct {
	code []string // TODO: [][]rune
	w, h int
	opts Opts

	done     bool // program has exited // TODO: move to exec
	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack
}

func NewProg(code string, opts Opts) *Prog {
	code = strings.TrimLeft(code, "\n")
	code = strings.TrimRightFunc(code, unicode.IsSpace)

	// TODO check for ASCII/Unicode input

	prog := Prog{
		code: strings.Split(code, "\n"),
		opts: opts,
	}

	prog.h = len(prog.code)
	for _, l := range prog.code {
		if len(l) > prog.w {
			prog.w = len(l)
		}
	}

	for i, l := range prog.code {
		// pad each row to width
		prog.code[i] = l + strings.Repeat(" ", prog.w-len(l))
	}

	return &prog
}

func (p *Prog) String() string {
	return strings.Join(p.code, "\n")
}

func (p *Prog) currentOp() opcode {
	return opcode(p.code[p.pcY][p.pcX])
}

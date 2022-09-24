package bef93

import (
	"io"
	"sync"
)

type Proc struct {
	prog Prog

	in          io.Reader
	out, outErr io.Writer

	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack
	done     bool

	//lint:ignore U1000 ignore unused copy guard
	noCopy sync.Mutex
}

func NewProc(prog *Prog, in io.Reader, out, outErr io.Writer) *Proc {
	return &Proc{
		prog: prog.Clone(),

		in:     in,
		out:    out,
		outErr: outErr,
	}
}

func (p *Proc) Clone() *Proc {
	return &Proc{
		prog: p.prog.Clone(),

		in:     p.in,
		out:    p.out,
		outErr: p.outErr,

		dir:     p.dir,
		pcX:     p.pcX,
		pcY:     p.pcY,
		strMode: p.strMode,
		stack:   p.stack.clone(),
		done:    p.done,
	}
}

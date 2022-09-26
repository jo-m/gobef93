package bef93

import (
	"io"
	"math/rand"
	"sync"
	"time"
)

type Proc struct {
	prog Prog

	in          io.Reader
	out, outErr io.Writer

	rand *rand.Rand

	dir      direction
	pcX, pcY int
	strMode  bool
	stack    stack
	done     bool

	//lint:ignore U1000 ignore unused copy guard
	noCopy sync.Mutex
}

func NewProc(prog *Prog, in io.Reader, out, outErr io.Writer) *Proc {
	seed := time.Now().UnixNano()
	if prog.opts.RandSeed != 0 {
		seed = prog.opts.RandSeed
	}

	return &Proc{
		prog: prog.Clone(),

		rand: rand.New(rand.NewSource(seed)),

		in:     in,
		out:    out,
		outErr: outErr,
	}
}

func (p *Proc) Prog() *Prog {
	prog := p.prog.Clone()
	return &prog
}

func (p *Proc) Clone(in io.Reader, out, outErr io.Writer) *Proc {
	return &Proc{
		prog: p.prog.Clone(),

		in:     in,
		out:    out,
		outErr: outErr,

		dir:     p.dir,
		pcX:     p.pcX,
		pcY:     p.pcY,
		strMode: p.strMode,
		stack:   p.stack.clone(),
		done:    p.done,
	}
}

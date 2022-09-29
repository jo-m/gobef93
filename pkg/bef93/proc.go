package bef93

import (
	"bufio"
	"io"
	"math/rand"
	"sync"
	"time"
)

// Proc represents a program in execution.
// Do not copy by value. Use proc.Clone() to obtain copies.
// Construct using NewProc().
type Proc struct {
	prog Prog

	in          *bufio.Reader
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

// NewProc creates a new Proc.
// In is the new procs stdin, out stdout, outErr stderr.
func NewProc(prog *Prog, in io.Reader, out, outErr io.Writer) *Proc {
	seed := time.Now().UnixNano()
	if prog.opts.RandSeed != 0 {
		seed = prog.opts.RandSeed
	}

	return &Proc{
		prog: prog.Clone(),

		// #nosec G404 We want to be deterministic here.
		rand: rand.New(rand.NewSource(seed)),

		in:     bufio.NewReader(in),
		out:    out,
		outErr: outErr,
	}
}

// Prog returns a copy of the current Prog inside the proc.
func (p *Proc) Prog() *Prog {
	prog := p.prog.Clone()
	return &prog
}

// Clone returns a pointer to a deep copy of a proc.
// You need to supply new I/O pipes.
func (p *Proc) Clone(in io.Reader, out, outErr io.Writer) *Proc {
	return &Proc{
		prog: p.prog.Clone(),

		in:     bufio.NewReader(in),
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

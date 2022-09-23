package bef93

import "sync"

type Proc struct {
	prog *Prog

	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack
	done     bool

	noCopy sync.Mutex
}

func NewProc(prog *Prog) *Proc {
	return &Proc{
		prog: prog.Clone(),
	}
}

func (p *Proc) Clone() *Proc {
	return &Proc{
		prog: p.prog.Clone(),

		dir:     p.dir,
		pcX:     p.pcX,
		pcY:     p.pcY,
		strMode: p.strMode,
		stack:   p.stack, // TODO: clone stack properly
		done:    p.done,
	}
}

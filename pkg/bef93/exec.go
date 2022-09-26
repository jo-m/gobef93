package bef93

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrTerminated    = errors.New("process already executed")
	ErrUnknownOpCode = errors.New("unknown opcode")
	ErrDivZero       = errors.New("division by zero")
	ErrWroteNothing  = errors.New("wrote 0 bytes")
)

var (
	// internal, not an actual error
	errTerminated = errors.New("process terminated")
)

func (p *Proc) Exec() error {
	if p.done {
		return ErrTerminated
	}
	defer func() { p.done = true }()

	for {
		err := p.step()
		if err == errTerminated {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (p *Proc) newRuntimeError(err error) *RuntimeError {
	return &RuntimeError{
		Msg:  err.Error(),
		Prog: p.prog.Clone(),
		LocX: p.pcX,
		LocY: p.pcY,

		cause: err,
	}
}

func (p *Proc) currentOp() opcode {
	return opcode(p.prog.code[p.pcY][p.pcX])
}

func (p *Proc) advancePC() {
	switch p.dir {
	case dirRight:
		p.pcX = (p.pcX + 1) % p.prog.w
	case dirDown:
		p.pcY = (p.pcY + 1) % p.prog.h
	case dirLeft:
		p.pcX = (p.pcX - 1 + p.prog.w) % p.prog.w
	case dirUp:
		p.pcY = (p.pcY - 1 + p.prog.h) % p.prog.h
	}
}

func (p *Proc) step() error {
	op := p.currentOp()
	iop := int64(op)

	if p.strMode && op != opStr {
		p.stack.push(iop)
	} else if strings.Contains("0123456789", string(op)) {
		p.stack.push(iop - int64('0'))
	} else {
		err := p.handleOp(op)
		if err != nil {
			return err
		}
	}

	p.advancePC()
	return nil
}

func readInt(in io.Reader) (int64, error) {
	r := bufio.NewReader(in)
	l, _, err := r.ReadLine()
	if err != nil {
		return 0, err
	}
	val, err := strconv.ParseInt(string(l), 10, 64)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func (p *Proc) handleOp(op opcode) error {
	switch op {
	case opAdd:
		a, b := p.stack.pop2()
		p.stack.push(a + b)
	case opSub:
		a, b := p.stack.pop2()
		p.stack.push(b - a)
	case opMul:
		a, b := p.stack.pop2()
		p.stack.push(a * b)
	case opDiv:
		a, b := p.stack.pop2()
		if a == 0 {
			if p.prog.opts.DisallowDivZero {
				return p.newRuntimeError(fmt.Errorf("%w: %d / %d", ErrDivZero, a, b))
			}

			fmt.Fprintf(p.outErr, "What do you want %d/0 to be?\n", b)
			b, err := readInt(p.in)
			if err != nil {
				if p.prog.opts.TerminateOnIOErr {
					return p.newRuntimeError(err)
				}

				b = 0
			}

			p.stack.push(b)
		} else {
			p.stack.push(b / a)
		}

	case opMod:
		a, b := p.stack.pop2()
		// in the reference implementation, this is not handled and would crash
		if a == 0 {
			return p.newRuntimeError(fmt.Errorf("%w: %d %% %d", ErrDivZero, a, b))
		}
		p.stack.push(b % a)
	case opNot:
		a := p.stack.pop()
		if a == 0 {
			p.stack.push(1)
		} else {
			p.stack.push(0)
		}
	case opGt:
		a, b := p.stack.pop2()
		if b > a {
			p.stack.push(1)
		} else {
			p.stack.push(0)
		}
	case opRight:
		p.dir = dirRight
	case opLeft:
		p.dir = dirLeft
	case opUp:
		p.dir = dirUp
	case opDown:
		p.dir = dirDown
	case opRand:
		p.dir = direction(p.rand.Intn(int(dirEND)))
	case opRif:
		a := p.stack.pop()
		if a == 0 {
			p.dir = dirRight
		} else {
			p.dir = dirLeft
		}
	case opDif:
		a := p.stack.pop()
		if a == 0 {
			p.dir = dirDown
		} else {
			p.dir = dirUp
		}
	case opStr:
		p.strMode = !p.strMode
	case opDup:
		a := p.stack.pop()
		p.stack.push(a)
		p.stack.push(a)
	case opSwp:
		a, b := p.stack.pop2()
		p.stack.push(a)
		p.stack.push(b)
	case opPop:
		_ = p.stack.pop()
	case opPopWrtInt:
		a := p.stack.pop()
		str := []byte(fmt.Sprintf("%d ", a))
		n, err := p.out.Write(str)
		if p.prog.opts.TerminateOnIOErr {
			if err != nil {
				return p.newRuntimeError(err)
			}
			if n == 0 {
				return p.newRuntimeError(ErrWroteNothing)
			}
		}
	case opPopWrtChr:
		chr := p.stack.pop()
		c := rune(chr)
		if !p.prog.opts.AllowUnicode {
			c = rune(byte(c))
		}
		n, err := p.out.Write([]byte(string([]rune{c})))
		if p.prog.opts.TerminateOnIOErr {
			if err != nil {
				return p.newRuntimeError(err)
			}
			if n == 0 {
				return p.newRuntimeError(ErrWroteNothing)
			}
		}
	case opSkip:
		p.advancePC()
	case opPut:
		y, x := p.stack.pop2()
		val := p.stack.pop()
		// TODO handle out of bounds
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		if !p.prog.opts.AllowUnicode {
			p.prog.code[y][x] = rune(byte(val))
		} else {
			p.prog.code[y][x] = rune(val)
		}
	case opGet:
		y, x := p.stack.pop2()
		// TODO handle out of bounds
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		val := p.prog.code[y][x]
		if !p.prog.opts.AllowUnicode {
			p.stack.push(int64(byte(val)))
		} else {
			p.stack.push(int64(val))
		}
	case opReadNr:
		val, err := readInt(p.in)
		if err != nil {
			if p.prog.opts.TerminateOnIOErr {
				return p.newRuntimeError(err)
			}

			if p.prog.opts.ReadErrorUndefined {
				// simulate "undefined" by using rand
				val = p.rand.Int63()
				if p.rand.Intn(2) == 0 {
					val = -val
				}
			} else {
				val = -1
			}
		}
		p.stack.push(val)
	case opReadChr:
		// TODO: this does not allow unicode
		buf := []byte{0}
		n, err := p.in.Read(buf)
		if err != nil && p.prog.opts.TerminateOnIOErr {
			return p.newRuntimeError(err)
		}
		if n == 0 || err != nil {
			// simulate EOF
			p.stack.push(-1)
		} else {
			p.stack.push(int64(buf[0]))
		}
	case opEnd:
		return errTerminated
	case opWhitespace:
		// do nothing
	default:
		if p.prog.opts.IgnoreUnsupportedInstructions {
			// do nothing
		} else {
			return p.newRuntimeError(fmt.Errorf("%w: '%s' (%d)", ErrUnknownOpCode, string(op), int64(op)))
		}
	}

	return nil
}

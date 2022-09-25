package bef93

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
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

	// log.Printf("new PC: %d, %d", p.pcX, p.pcY)
}

func (p *Proc) step() error {
	op := p.currentOp()
	iop := int64(op)

	if p.strMode && op != opStr {
		log.Println("str push", op, iop)
		p.stack.push(iop)
	} else if strings.Contains("0123456789", string(op)) {
		log.Println("num push", iop-int64('0'))
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
		log.Println("opAdd")
		a, b := p.stack.pop2()
		p.stack.push(a + b)
		log.Printf("opAdd %d+%d=%d", a, b, a+b)
	case opSub:
		log.Println("opSub")
		a, b := p.stack.pop2()
		p.stack.push(b - a)
		log.Printf("opSub %d-%d=%d", b, a, b-a)
	case opMul:
		log.Println("opMul")
		a, b := p.stack.pop2()
		p.stack.push(a * b)
		log.Printf("opMul %d*%d=%d", a, b, a*b)
	case opDiv:
		log.Println("opDiv")
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

		log.Printf("opDiv %d/%d=%d", b, a, b/a)
	case opMod:
		log.Println("opMod")
		a, b := p.stack.pop2()
		// this is not handled in standard Befunge93
		if a == 0 {
			return p.newRuntimeError(fmt.Errorf("%w: %d %% %d", ErrDivZero, a, b))
		}
		p.stack.push(b % a)
		log.Printf("opMod %d%%%d=%d", b, a, b%a)
	case opNot:
		log.Println("opNot")
		a := p.stack.pop()
		if a == 0 {
			p.stack.push(1)
		} else {
			p.stack.push(0)
		}
		log.Printf("opNot !%d", a)
	case opGt:
		log.Println("opGt")
		a, b := p.stack.pop2()
		if b > a {
			p.stack.push(1)
		} else {
			p.stack.push(0)
		}
		log.Printf("opDiv %d>%d", b, a)
	case opRight:
		log.Println("dirRight")
		p.dir = dirRight
	case opLeft:
		log.Println("dirLeft")
		p.dir = dirLeft
	case opUp:
		log.Println("dirUp")
		p.dir = dirUp
	case opDown:
		log.Println("dirDown")
		p.dir = dirDown
	case opRand:
		p.dir = direction(p.rand.Intn(int(dirEND)))
		log.Println("opRand", p.dir)
	case opRif:
		log.Println("opRif")
		a := p.stack.pop()
		if a == 0 {
			p.dir = dirRight
		} else {
			p.dir = dirLeft
		}
		log.Println("opRif", p.dir)
	case opDif:
		log.Println("opDif")
		a := p.stack.pop()
		if a == 0 {
			p.dir = dirDown
		} else {
			p.dir = dirUp
		}
		log.Println("opDif", p.dir)
	case opStr:
		p.strMode = !p.strMode
		log.Printf("string mode enabled: %t", p.strMode)
	case opDup:
		log.Println("opDup")
		a := p.stack.pop()
		p.stack.push(a)
		p.stack.push(a)
		log.Println("opDup", a)
	case opSwp:
		log.Println("opSwp")
		a, b := p.stack.pop2()
		p.stack.push(a)
		p.stack.push(b)
		log.Println("opSwp", a, b)
	case opPop:
		log.Println("opPop")
		_ = p.stack.pop()
	case opPopWrtInt:
		log.Println("opPopWrtInt")
		a := p.stack.pop()
		str := []byte(fmt.Sprint(a))
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
		log.Println("opPopWrtChr")
		chr := p.stack.pop()
		c := byte(chr) // TODO unicode support
		n, err := p.out.Write([]byte{c})
		if p.prog.opts.TerminateOnIOErr {
			if err != nil {
				return p.newRuntimeError(err)
			}
			if n == 0 {
				return p.newRuntimeError(ErrWroteNothing)
			}
		}
	case opSkip:
		log.Println("opSkip")
		p.advancePC()
	case opPut:
		log.Println("opPut")
		y, x := p.stack.pop2()
		val := p.stack.pop()
		// TODO how to handle out of bounds
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		// TODO: handle unicode
		p.prog.code[y][x] = rune(val)
		log.Println("opPut", x, y, val)
	case opGet:
		log.Println("opGet")
		y, x := p.stack.pop2()
		// TODO how to handle out of bounds
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		val := p.prog.code[y][x]
		// TODO: handle unicode
		p.stack.push(int64(val))
		log.Println("opGet", x, y, val)
	case opReadNr:
		log.Println("opReadNr")
		val, err := readInt(p.in)
		if err != nil && p.prog.opts.TerminateOnIOErr {
			return p.newRuntimeError(err)
		}
		if err != nil {
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
		log.Println("opReadNr", int(val))
	case opReadChr:
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
		return p.newRuntimeError(fmt.Errorf("%w: '%s' (%d)", ErrUnknownOpCode, string(op), int64(op)))
	}

	return nil
}

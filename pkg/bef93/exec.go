package bef93

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

var (
	ErrTerminated    = errors.New("process already executed")
	ErrUnknownOpCode = errors.New("unknown opcode")
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

// TODO make usage consistent and sensible
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
		// TODO div by 0
		p.stack.push(b / a)
		log.Printf("opDiv %d/%d=%d", b, a, b/a)
	case opMod:
		log.Println("opMod")
		a, b := p.stack.pop2()
		// TODO div by 0
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
		// TODO: allow to seed rand
		p.dir = direction(rand.Intn(int(dirEND)))
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
		// TODO: how to do io error handling
		if n != len(str) {
			return fmt.Errorf("failed to write: %w", err)
		}
	case opPopWrtChr:
		log.Println("opPopWrtChr")
		chr := p.stack.pop()
		if chr < 0 || chr > math.MaxUint8 {
			return errors.New("overflow")
		}
		c := byte(chr)
		n, err := p.out.Write([]byte{c})
		// TODO: how to do io error handling
		if n != 1 {
			return fmt.Errorf("failed to write: %w", err)
		}
	case opSkip:
		log.Println("opSkip")
		p.advancePC()
	case opPut:
		log.Println("opPut")
		y, x := p.stack.pop2()
		val := p.stack.pop()
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		// TODO: how to handle overflow
		p.prog.code[y][x] = rune(val)
		log.Println("opPut", x, y, val)
	case opGet:
		log.Println("opGet")
		y, x := p.stack.pop2()
		y = (y + int64(p.prog.h)) % int64(p.prog.h)
		x = (x + int64(p.prog.w)) % int64(p.prog.w)
		val := p.prog.code[y][x]
		// TODO: how to handle overflow
		p.stack.push(int64(val))
		log.Println("opGet", x, y, val)
	case opReadNr:
		log.Println("opReadNr")
		// TODO: how to handle overflow
		_, err := fmt.Fprintln(p.outErr, "> Enter an integer and press Enter:")
		if err != nil {
			return err
		}
		r := bufio.NewReader(p.in)
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		val, err := strconv.ParseInt(string(l), 10, 32)
		if err != nil {
			return err
		}
		p.stack.push(int64(val))
		log.Println("opReadNr", int(val))
	case opReadChr:
		// TODO: implement
		// TODO: how to handle overflow
		panic("opReadChr not implemented")
	case opEnd:
		return errTerminated
	case opWhitespace:
		// do nothing
	default:
		return p.newRuntimeError(fmt.Errorf("code '%s' (%d): %w", string(op), int64(op), ErrUnknownOpCode))
		// return fmt.Errorf("unknown opcode '%c' at (%d,%d)", op, p.pcX, p.pcY)
	}

	return nil
}

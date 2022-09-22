package befunge

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type Prog struct {
	code []string
	w, h int // size

	done     bool // program has exited
	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack
}

func NewProg(code string) *Prog {
	code = strings.TrimSpace(code)

	prog := Prog{
		code: strings.Split(code, "\n"),
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

func (p *Prog) advancePC() {
	switch p.dir {
	case dirRight:
		p.pcX = (p.pcX + 1) % p.w
	case dirDown:
		p.pcY = (p.pcY + 1) % p.h
	case dirLeft:
		p.pcX = (p.pcX - 1 + p.w) % p.w
	case dirUp:
		p.pcY = (p.pcY - 1 + p.h) % p.h
	}

	// log.Printf("new PC: %d, %d", p.pcX, p.pcY)
}

func (p *Prog) handleOp(op opcode, in io.Reader, out, outErr io.Writer) error {
	switch op {
	case opAdd:
		log.Println("opAdd")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		p.stack.push(a + b)
		log.Printf("opAdd %d+%d=%d", a, b, a+b)
	case opSub:
		log.Println("opSub")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		p.stack.push(b - a)
		log.Printf("opSub %d-%d=%d", b, a, b-a)
	case opMul:
		log.Println("opMul")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		p.stack.push(a * b)
		log.Printf("opMul %d*%d=%d", a, b, a*b)
	case opDiv:
		log.Println("opDiv")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		// TODO div by 0
		p.stack.push(b / a)
		log.Printf("opDiv %d/%d=%d", b, a, b/a)
	case opMod:
		log.Println("opMod")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		p.stack.push(b % a)
		log.Printf("opMod %d%%%d=%d", b, a, b%a)
	case opNot:
		log.Println("opNot")
		a, err := p.stack.pop()
		if err != nil {
			return err
		}
		if a == 0 {
			p.stack.push(1)
		} else {
			p.stack.push(0)
		}
		log.Printf("opNot !%d", a)
	case opGt:
		log.Println("opGt")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
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
		p.dir = direction(rand.Intn(int(dirEND)))
		log.Println("opRand", p.dir)
	case opRif:
		log.Println("opRif")
		a, err := p.stack.pop()
		if err != nil {
			return err
		}
		if a == 0 {
			p.dir = dirRight
		} else {
			p.dir = dirLeft
		}
		log.Println("opRif", p.dir)
	case opDif:
		log.Println("opDif")
		a, err := p.stack.pop()
		if err != nil {
			return err
		}
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
		a, err := p.stack.pop()
		if err != nil {
			return err
		}
		p.stack.push(a)
		p.stack.push(a)
		log.Println("opDup", a)
	case opSwp:
		log.Println("opSwp")
		a, b, err := p.stack.pop2()
		if err != nil {
			return err
		}
		p.stack.push(a)
		p.stack.push(b)
		log.Println("opSwp", a, b)
	case opPop:
		log.Println("opPop")
		_, err := p.stack.pop()
		if err != nil {
			return err
		}
	case opPopWrtInt:
		log.Println("opPopWrtInt")
		a, err := p.stack.pop()
		if err != nil {
			return err
		}
		str := []byte(fmt.Sprint(a))
		n, err := out.Write(str)
		if n != len(str) {
			return errors.New("failed to write")
		}
		if err != nil {
			return err
		}
	case opPopWrtChr:
		log.Println("opPopWrtChr")
		chr, err := p.stack.pop()
		if err != nil {
			return err
		}
		if chr < 0 || chr > math.MaxUint8 {
			return errors.New("overflow")
		}
		c := byte(chr)
		n, err := out.Write([]byte{c})
		if n != 1 {
			return errors.New("failed to write")
		}
		if err != nil {
			return err
		}
	case opSkip:
		log.Println("opSkip")
		p.advancePC()
	case opPut:
		log.Println("opPut")
		y, x, err := p.stack.pop2()
		if err != nil {
			return err
		}
		val, err := p.stack.pop()
		if err != nil {
			return err
		}
		y = (y + p.h) % p.h
		x = (x + p.w) % p.w
		str := []byte(p.code[y])
		str[x] = byte(val)
		p.code[y] = string(str)
		log.Println("opPut", x, y, val)
	case opGet:
		log.Println("opGet")
		y, x, err := p.stack.pop2()
		if err != nil {
			return err
		}
		y = (y + p.h) % p.h
		x = (x + p.w) % p.w
		val := p.code[y][x]
		p.stack.push(int(val))
		log.Println("opGet", x, y, val)
	case opReadNr:
		log.Println("opReadNr")
		_, err := fmt.Fprintln(outErr, "> Enter an integer and press Enter:")
		if err != nil {
			return err
		}
		r := bufio.NewReader(in)
		l, _, err := r.ReadLine()
		if err != nil {
			return err
		}
		val, err := strconv.ParseInt(string(l), 10, 32)
		if err != nil {
			return err
		}
		p.stack.push(int(val))
		log.Println("opReadNr", int(val))
	case opReadChr:
		panic("opReadChr not implemented")
	case opEnd:
		panic("should be handled in main loop")
	case opWhitespace:
		// do nothing
	default:
		// return fmt.Errorf("unknown opcode '%c' at (%d,%d)", op, p.pcX, p.pcY)
	}

	return nil
}

func (p *Prog) Exec(in io.Reader, out, outErr io.Writer) error {
	if p.done {
		return errors.New("already executed")
	}
	defer func() { p.done = true }()

	for i := 0; i < 100000000; i++ { // TODO infinite
		op := p.currentOp()
		iop := int(op)

		if p.strMode && op != opStr {
			// string mode
			log.Println("str push", op, iop)
			p.stack.push(iop)
		} else if iop >= int('0') && iop <= int('9') {
			// numbers
			log.Println("num push", iop-int('0'))
			p.stack.push(iop - int('0'))
		} else if op == opEnd {
			return nil
		} else {
			err := p.handleOp(op, in, out, outErr)
			if err != nil {
				return err
			}
		}

		p.advancePC()
	}

	panic("did not terminate")
}

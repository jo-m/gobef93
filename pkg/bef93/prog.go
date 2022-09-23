package bef93

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
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
	AllowArbitraryCodeSize bool // allow code of arbitrary size, code smaller than standard size will still be padded to standard size
	AllowUnicode           bool // TODO: implement
}

// Prog represents a Befunge-93 program.
// Use NewProg() to get an instance.
// Immutable after creation.
type Prog struct {
	code [][]rune // TODO: deep-copy for execution
	w, h int
	opts Opts

	// TODO move state to exec struct, to make this immutable
	done     bool // program has exited // TODO: move to exec
	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack

	noCopy sync.Mutex
}

type State struct {
	dir      direction
	pcX, pcY int  // program counter
	strMode  bool // string mode active
	stack    stack
}

type Proc struct {
	prog  Prog
	state State
	done  bool
}

func isASCII(lines []string) (bool, int, int) {
	for y, l := range lines {
		for x, c := range l {
			if c > unicode.MaxASCII {
				return false, x, y
			}
		}
	}

	return true, 0, 0
}

// returns max size, and its location
func getMaxSize(lines []string) (w, h int) {
	h = len(lines)

	for _, line := range lines {
		l := len(line)
		if l > w {
			w = l
		}
	}

	return
}

var (
	ErrNotASCII = errors.New("code contains non-ascii characters")
	ErrTooLarge = errors.New("program code is too large")
)

func NewProg(code string, opts Opts) (*Prog, error) {
	lines := strings.Split(code, "\n")

	if !opts.AllowUnicode {
		ok, x, y := isASCII(lines)
		if !ok {
			return nil, NewCompilationError(ErrNotASCII, x, y)
		}
	}

	w, h := getMaxSize(lines)
	if !opts.AllowArbitraryCodeSize && (w > Width || h > Height) {
		return nil, NewCompilationError(ErrTooLarge, w, h)
	}

	if !opts.AllowArbitraryCodeSize && (w > Width || h > Height) {
		return nil, NewCompilationError(ErrTooLarge, w, h)
	}
	if w < Width {
		w = Width
	}
	if h < Height {
		h = Height
	}

	// pad rows
	for len(lines) < h {
		lines = append(lines, "")
	}

	// pad cols
	for i, l := range lines {
		lines[i] = l + strings.Repeat(" ", w-len(l))
	}

	runes := make([][]rune, h)
	for i, l := range lines {
		runes[i] = []rune(l)
	}

	return &Prog{
		code: runes,
		w:    w,
		h:    h,
		opts: opts,
	}, nil
}

func (p *Prog) String() string {
	b := strings.Builder{}

	numSz := int(math.Ceil(math.Log10(math.Max(float64(p.w), float64(p.h))))) + 1

	// top numbering
	b.WriteString("    ")
	for x := 0; x < p.w; x += 10 {
		b.WriteString(fmt.Sprintf("v-%-8d", x))
	}
	b.WriteString("\n")

	// code and side numbering
	b.WriteString(strings.Repeat(" ", numSz) + "|" + strings.Repeat("-", p.w) + "|\n")
	for y, l := range p.code {
		b.WriteString(fmt.Sprintf("% *d|", numSz, y))
		for _, r := range l {
			b.WriteRune(r)
		}
		b.WriteString(fmt.Sprintf("|% *d\n", numSz, y))
	}
	b.WriteString(strings.Repeat(" ", numSz) + "|" + strings.Repeat("-", p.w) + "|\n")

	// bottom numbering
	b.WriteString("    ")
	for x := 0; x < p.w; x += 10 {
		b.WriteString(fmt.Sprintf("^-%-8d", x))
	}

	return b.String()
}

package bef93

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"unicode"
)

// default program size
const (
	Width  = 80
	Height = 25
)

// Opts contains supported options.
// See https://github.com/catseye/Befunge-93/blob/master/src/bef.c#L46.
// Zero value is good to use and represents the default (standard 93) options.
// If you update docstrings and options here, also update them in main.go.
type Opts struct {
	// Options (mostly) equal to reference implementation.

	// TODO: implement.
	NoFixOffByOne bool
	// If true, & will push an undefined number to stack instead of -1.
	ReadErrorUndefined bool
	// If true, unsupported instructions will be ignored.
	// Note that unline the reference implementation, we
	// terminate instead of just print an error on unsupported
	// instructions by default.
	IgnoreUnsupportedInstructions bool
	// TODO: implement.
	WrapLongLines bool
	// TODO: implement.
	WrapHashInconsistently bool

	// Non-standard options.

	// Allow code of arbitrary size, code smaller
	// than standard size will be padded to standard size.
	AllowArbitraryCodeSize bool
	// Allow unicode in the interpreted code.
	// This also allows the 'g' and 'p' operators to load/store unicode runes,
	// and the ',' and '~' operators to write/read unicode runes (utf-8 encoded).
	AllowUnicode bool
	// Terminate on division by 0.
	DisallowDivZero bool
	// Fixed random seed. If 0, the generator
	// is seeded randomly internally.
	// This allows to deterministically execute programs containing
	// random operations.
	RandSeed int64
	// Terminate on I/O errors instead of ignoring them.
	TerminateOnIOErr bool
	// Terminate if a 'g' or 'p' operation is out of bounds, instead of pushing 0 or discading the pop() value.
	TerminateOnPutGetOutOfBounds bool
}

// Prog represents a Befunge-93 program.
// Use NewProg() to get an instance.
// Do not copy by value, use prog.Clone() to obtain copies.
type Prog struct {
	code [][]rune
	w, h int
	opts Opts

	//lint:ignore U1000 ignore unused copy guard.
	// Do not create naive struct copies, use p.Clone() instead.
	noCopy sync.Mutex
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

// Common errors returned by NewProg().
// Will be wrapped in a CompilationError, so use errors.Is/As().
var (
	ErrNotASCII = errors.New("code contains non-ascii characters")
	ErrTooLarge = errors.New("program code is too large")
)

// NewProg creates a new program from source code and options.
func NewProg(code string, opts Opts) (*Prog, error) {
	lines := strings.Split(code, "\n")

	if !opts.AllowUnicode {
		ok, x, y := isASCII(lines)
		if !ok {
			return nil, newCompilationError(ErrNotASCII, x, y)
		}
	}

	w, h := getMaxSize(lines)
	if !opts.AllowArbitraryCodeSize && (w > Width || h > Height) {
		return nil, newCompilationError(ErrTooLarge, w, h)
	}

	if !opts.AllowArbitraryCodeSize && (w > Width || h > Height) {
		return nil, newCompilationError(ErrTooLarge, w, h)
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

	if opts.NoFixOffByOne {
		panic("option NoFixOffByOne: not implemented")
	}
	if opts.WrapLongLines {
		panic("option WrapLongLines: not implemented")
	}
	if opts.WrapHashInconsistently {
		panic("option WrapHashInconsistently: not implemented")
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

// Code returns the source code of this program.
func (p *Prog) Code() string {
	ret := strings.Builder{}
	for _, l := range p.code {
		ret.WriteString(strings.TrimRight(string(l), " "))
		ret.WriteByte('\n')
	}
	return strings.TrimSpace(ret.String())
}

// Opts returns the options of this program.
func (p *Prog) Opts() Opts {
	return p.opts
}

// Clone returns a pointer to a deep copy of a prog.
func (p *Prog) Clone() Prog {
	code := make([][]rune, p.h)
	for i, line := range p.code {
		code[i] = make([]rune, len(line))
		copy(code[i], line)
	}

	return Prog{
		code: code,
		w:    p.w,
		h:    p.h,
		opts: p.opts,
	}
}

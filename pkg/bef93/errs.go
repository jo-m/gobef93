package bef93

type CompilationError struct {
	Msg        string // error message
	LocX, LocY int    // error location in code

	cause error
}

// compile time interface check
var _ error = (*CompilationError)(nil)

func NewCompilationError(err error, locX, locY int) *CompilationError {
	return &CompilationError{
		Msg:  err.Error(),
		LocX: locX,
		LocY: locY,

		cause: err,
	}
}

func (e *CompilationError) Error() string { return e.Msg }

func (e *CompilationError) Unwrap() error { return e.cause }

type RuntimeError struct {
	Msg        string
	Prog       Prog
	LocX, LocY int

	cause error
}

// compile time interface check
var _ error = (*RuntimeError)(nil)

func (e *RuntimeError) Error() string { return e.Msg }

func (e *RuntimeError) Unwrap() error { return e.cause }

// Backtrace returns a user-friendly error message describing the state of the program
// that led to this error.
func (e *RuntimeError) Backtrace() string {
	panic("not implemented") // TODO: implement
}

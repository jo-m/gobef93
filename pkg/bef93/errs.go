package bef93

import "fmt"

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

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error at (%d, %d): %s", e.LocX, e.LocY, e.Msg)
}

func (e *RuntimeError) Unwrap() error { return e.cause }

package errors

import (
	"fmt"
	"log"
)

// Delim and UnitSep are used to separate/delimit fields
const (
	Delim   = "\n"
	UnitSep = ";\037" // []byte{31} is the ascii Unit Separator (US) character
)

// Error is an error with an embedded "previous" error
type Error struct {
	err  string
	prev error
}

// New creates a new Error of our own liking
func New(args ...interface{}) error {
	if len(args) == 0 {
		log.Fatal("call to errors.New with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.err = string(arg)
		case *Error:
			e.prev = arg
		case error:
			e.prev = New(arg.Error)
		}
	}
	return e
}

// Prev allows access to an Error's prev
func (e *Error) Prev() error {
	return e.prev
}

// Error fulfills the error interface
func (e *Error) Error() string {
	return string(Serialize(e, nil, Delim))
}

// Error fulfills the stringer interface
func (e *Error) String() string {
	return string(Serialize(e, nil, "; "))
}

// Serialize writes the entire stack using sep as a delimeter
func (e *Error) Serialize(b []byte, sep string) []byte {
	b = append(b, e.err...)
	if e.prev != nil {
		b = append(b, sep...)
		b = Serialize(e.prev, b, sep)
	}
	return b
}

// Serialize writes the entire stack using sep as a delimeter
func Serialize(err error, b []byte, sep string) []byte {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e.Serialize(b, sep)
	}
	b = append(b, err.Error()...)
	return b
}

// Pulled from https://godoc.org/upspin.io/errors

// Recreate the errors.New functionality of the standard Go errors package
// so we can create simple text errors when needed.

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}

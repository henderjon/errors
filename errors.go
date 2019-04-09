package errors

import (
	"bytes"
	"fmt"
	"log"
)

const delim = "\n"
const unitsep = ";\037" // []byte{31} is the ascii Unit Separator (US) character

// Error is an error with an embedded "previous" error
type Error struct {
	err  string
	prev *Error
}

// New creates a new Error of our own liking
func New(args ...interface{}) *Error {
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
func (e *Error) Prev() *Error {
	return e.prev
}

// Error fulfills the error interface
func (e *Error) Error() string {
	return string(e.Marshal(delim))
}

// Error fulfills the stringer interface
func (e *Error) String() string {
	return string(e.Marshal("; "))
}

// Marshal writes the entire stack using sep as a delimeter
func (e *Error) Marshal(sep string) []byte {
	var b bytes.Buffer
	b.WriteString(e.err)
	if e.prev != nil {
		b.WriteString(sep)
		b.Write(e.prev.Marshal(sep))
	}
	return b.Bytes()
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling. Pulled from https://godoc.org/upspin.io/errors
func Errorf(format string, args ...interface{}) error {
	return New(fmt.Sprintf(format, args...))
}

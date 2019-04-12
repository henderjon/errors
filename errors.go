package errors

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

// Delim and UnitSep are used to separate/delimit fields
const (
	Delim   = "\n\t"
	UnitSep = ";\037" // []byte{31} is the ascii Unit Separator (US) character
)

// Error is an error with an embedded "previous" error
type Error struct {
	err  string
	prev error
	loc  Location
}

// New creates a new Error of our own liking. The args passed should be the
// the current error string and the previous error as either a standard error or
// an *Error from this package.
func New(args ...interface{}) error {
	if len(args) == 0 {
		log.Fatal("call to errors.New with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.err = arg
		case Location:
			e.loc = arg
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
	if e.loc != "" {
		b = append(b, " @ "...)
		b = append(b, e.loc...)
	}
	if e.prev != nil {
		b = append(b, sep...)
		b = Serialize(e.prev, b, sep)
	}
	return b
}

// Serialize writes the entire stack using sep as a delimeter. The args passed
// should be the latest (topmost) error and either a delim/separator string and
// a []byte to populate. The []byte arg will see almost no usage as it's
// primarily used for the recursive stack building. Although it's certainly not
// out of the realm of possibility that there is a []byte to be filled.
func Serialize(err error, args ...interface{}) []byte {
	if err == nil {
		return nil
	}

	var (
		b   []byte
		sep = Delim
	)

	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			sep = arg
		case []byte:
			b = arg
		}
	}

	if e, ok := err.(*Error); ok {
		return e.Serialize(b, sep)
	}

	b = append(b, err.Error()...)
	return b
}

// Location is the name:line of a file. Ideally returned by Here(). In usage
// it'll give you the file:line of the invocation of Here() to be passed as part
// of the error.
type Location string

// Here return the file:line of calling Here()
func Here() Location {
	var l Location
	_, file, line, ok := runtime.Caller(1)
	if ok {
		path := filepath.Base(file)
		l = Location(path + ":" + strconv.Itoa(line))
	}
	return l
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

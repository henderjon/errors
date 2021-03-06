package errors

import (
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// Sep and UnitSep are used to separate/delim fields
const (
	Sep       = "\n\t"
	RecordSep = "\036" // byte(30) || "\x1e" ... is the ascii Record Separator (RS) character
	UnitSep   = "\037" // byte(31) || "\x1e" ... is the ascii Unit Separator (US) character
)

// Error is an error with an embedded "previous" error and a kind
type Error struct {
	Err      string   `json:"error,omitempty"`    // this error
	Kind     Kind     `json:"kind,omitempty"`     // the Kind of this error
	Location Location `json:"location,omitempty"` // the location of this error
	Prev     error    `json:"previous,omitempty"` // the previous error
}

// Kind is a custom int type to communicate the error's Kind.
type Kind uint8

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

// New creates a new Error of our own liking. The `string` args are assumed
// to be the error message. The `error`/`Error` arg is assumed to be a Prev.
// The `Location` arg is assumed to be the Location. The `Kind` arg is the
// Kind of the error.
func New(args ...interface{}) error {
	if len(args) == 0 {
		log.Fatal("call to errors.New with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Err = arg
		case Kind:
			e.Kind = arg
		case Location:
			e.Location = arg
		case *Error:
			if len(args) == 1 {
				return arg // don't just wrap errors in errors for funsies
			}
			e.Prev = arg
		case error:
			if len(args) == 1 {
				return arg // don't just wrap errors in errors for funsies
			}
			e.Prev = New(arg.Error()) // I don't see any good reason to preserve the Stdlib's errors.
		}
	}
	return e
}

// Is checks to see if the error is of a certain kind
func Is(err error, k Kind) bool {
	if e, ok := err.(*Error); ok {
		return e.Kind == k
	}
	return false
}

// Has checks to see if the error is of a certain kind. It returns the (matching error, true) or the (nil, false)
func Has(err error, k Kind) (error, bool) {
	if e, ok := err.(*Error); ok {
		if e.Kind == k {
			return e, true
		}
		if e.Prev != nil {
			return Has(e.Prev, k)
		}
	}
	return nil, false
}

// Error fulfills the error interface. The error stack will be of the format:
// `@ location; message\n\t`
func (e *Error) Error() string {
	var b strings.Builder
	if e.Location != "" {
		b.WriteString("@ ")
		b.WriteString(string(e.Location))
		b.WriteString("; ")
	}
	if e.Err != "" {
		b.WriteString(e.Err)
	}
	if e.Prev != nil {
		b.WriteString(Sep)
		b.WriteString(e.Prev.Error())
	}
	return b.String()
}

// Serialize writes the entire stack using the format 'int64(len)[]bytes(value)'
func (e *Error) Serialize(b []byte) []byte {
	b = appendKind(b, e.Kind)
	b = appendString(b, string(e.Location))
	b = appendString(b, e.Err)
	b = appendError(b, e.Prev)
	return b
}

// appendString writes a string's length and value to b
func appendString(b []byte, str string) []byte {
	var tmp [16]byte // For use by PutUvarint.
	N := binary.PutUvarint(tmp[:], uint64(len(str)))
	b = append(b, tmp[:N]...)
	b = append(b, str...)
	return b
}

// appendString writes a Kind's length and value to b
func appendKind(b []byte, k Kind) []byte {
	var tmp [16]byte                        // For use by PutVarint.
	N := binary.PutVarint(tmp[:], int64(k)) // value
	b = append(b, tmp[:N]...)               // value
	return b
}

// appendError writes an Error/error to b
func appendError(b []byte, e error) []byte {
	if e == nil {
		return b
	}
	if e, ok := e.(*Error); ok {
		return e.Serialize(b)
	}
	return appendString(b, e.Error())
}

// Serialize writes the entire stack using a binary encoding. The args passed
// should be the latest (topmost) error and a []byte to populate.
// The []byte arg will see almost no usage as it's
// primarily used for the recursive serializing. Although it's certainly not
// out of the realm of possibility that there is a []byte to be filled.
func Serialize(err error, args ...interface{}) []byte {
	if err == nil {
		return nil
	}

	var b []byte

	for _, arg := range args {
		switch arg := arg.(type) {
		case []byte:
			b = arg // if a []byte was passed, fill it
		}
	}

	if e, ok := err.(*Error); ok {
		return e.Serialize(b)
	}

	b = appendString(b, err.Error())
	return b
}

// Unserialize reads the byte slice into the receiver, which must be non-nil.
// The returned error is always nil.
func (e *Error) Unserialize(b []byte) error {
	var data []byte
	if len(b) == 0 {
		return nil
	}

	// the first value is a scalar int64 `kind`
	k, N := binary.Varint(b)
	e.Kind = Kind(k)
	b = b[N:]

	data, b = getBytes(b)
	if data != nil {
		e.Location = Location(data)
	}

	data, b = getBytes(b)
	if data != nil {
		e.Err = string(data)
	}

	e.Prev = parseError(b)

	return nil
}

// parseError either unserializes the error or returns nil as no error is present
func parseError(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	var err Error
	err.Unserialize(b)
	return &err
}

// getBytes reads the byte slice at b (uvarint count followed by bytes)
// and returns the slice followed by the remaining bytes.
// If there is insufficient data, both return values will be nil.
func getBytes(b []byte) (data, remaining []byte) {
	u, N := binary.Uvarint(b)
	if len(b) < N+int(u) {
		log.Println("Unmarshal error: bad encoding 1")
		return nil, nil
	}
	if N == 0 {
		log.Println("Unmarshal error: bad encoding 2")
		return nil, b
	}
	return b[N : N+int(u)], b[N+int(u):]
}

// Encode takes the error and DSV encodes is for storage. It's a bit esoteric to use
// byte(30) and byte(31) as delimeters, but that's why those characters exist. Without
// type info, this isn't all that useful.
func Encode(e error) string {
	if e == nil {
		return ""
	}
	var b strings.Builder
	if e, ok := e.(*Error); ok {
		b.WriteString(fmt.Sprintf("%03d", e.Kind))
		b.WriteString(UnitSep)
		b.WriteString(string(e.Location))
		b.WriteString(UnitSep)
		b.WriteString(e.Err)
		b.WriteString(RecordSep)
		b.WriteString(Encode(e.Prev))
	} else {
		b.WriteString(e.Error())
	}

	return b.String()
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

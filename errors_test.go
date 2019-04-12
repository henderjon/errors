package errors

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestErrors(t *testing.T) {
	errs := []string{
		"first",
		"second",
		"third",
		"fourth",
	}

	tmpE := New("OH NOES!")
	for _, e := range errs {
		tmpE = New(e, tmpE)
	}

	if diff := cmp.Diff(tmpE.Error(), string(Serialize(tmpE, Delim))); diff != "" {
		t.Error("Error and Marshal should be of same len(); (-got +want)", diff)
	}

	expected := "fourth; third; second; first; OH NOES!"
	if diff := cmp.Diff(expected, string(Serialize(tmpE, "; "))); diff != "" {
		t.Error("unexpected Serialize(); (-got +want)", diff)
	}

	// fmt.Println(tmpE)
	// fmt.Println(Here())

	// fmt.Println(string(Serialize(tmpE, UnitSep)))

	// var b []byte
	// b = Serialize(tmpE, b, ";")
	// fmt.Println(string(b))

}

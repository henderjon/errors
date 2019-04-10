package errors

import (
	"testing"
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

	if tmpE.Error() != string(Serialize(tmpE, Delim)) {
		t.Fatal("Error and Marshal should be of same len()")
	}

	if "fourth; third; second; first; OH NOES!" != string(Serialize(tmpE, "; ")) {
		t.Fatal("Error and Marshal should be of same len()")
	}

	// fmt.Println(tmpE)

	// fmt.Println(string(Serialize(tmpE, UnitSep)))

	// var b []byte
	// b = Serialize(tmpE, b, ";")
	// fmt.Println(string(b))

}

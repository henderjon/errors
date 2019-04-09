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

	tmpE := New("")
	for _, e := range errs {
		tmpE = New(e, tmpE)
	}

	if len(tmpE.Error()) != len(tmpE.Marshal("\t")) {
		t.Fatal("Error and Marshal should be of same len()")
	}

}

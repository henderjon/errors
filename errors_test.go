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

	if len(tmpE.Error()) != len(Serialize(tmpE, nil, "-")) {
		t.Fatal("Error and Marshal should be of same len()")
	}

	// fmt.Println(tmpE)

}

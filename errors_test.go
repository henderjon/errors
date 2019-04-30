package errors

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	bad Kind = iota + 1
	worse
	worst
)

func getErrorForSerialization() error {
	a := New(bad, Here(), "things are gonna be bad")
	b := New(worse, "getErrorForSerialization", Here(), a)
	c := New(worst, Here(), b)
	return c
}

func TestEncode(t *testing.T) {
	e := getErrorForSerialization()

	expected := "003\x1ferrors_test.go:19\x1f\x1e002\x1ferrors_test.go:18\x1fgetErrorForSerialization\x1e001\x1ferrors_test.go:17\x1fthings are gonna be bad\x1e"
	if diff := cmp.Diff(Encode(e), expected); diff != "" {
		t.Error("Encode(e); (-got +want)", diff)
	}
}

func TestString(t *testing.T) {
	e := getErrorForSerialization()

	expected := `@ errors_test.go:19; ` + `
	@ errors_test.go:18; getErrorForSerialization
	@ errors_test.go:17; things are gonna be bad`
	if diff := cmp.Diff(e.Error(), expected); diff != "" {
		// t.Fatal(e.Error())
		t.Error("Error.Error(); (-got +want)", diff)
	}
}

func TestJSONEncode(t *testing.T) {
	e := getErrorForSerialization()

	real, _ := json.Marshal(e)
	// t.Fatal(string(real))
	expected := []byte(`{"kind":3,"location":"errors_test.go:19","previous":{"error":"getErrorForSerialization","kind":2,"location":"errors_test.go:18","previous":{"error":"things are gonna be bad","kind":1,"location":"errors_test.go:17"}}}`)
	if diff := cmp.Diff(real, expected); diff != "" {
		t.Error("json.Marshal(); (-got +want)", diff)
	}
}

func TestSerialize(t *testing.T) {
	e := getErrorForSerialization()

	expected := []byte{6, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 57, 0, 4, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 56, 24, 103, 101, 116, 69, 114, 114, 111, 114, 70, 111, 114, 83, 101, 114, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110, 2, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 55, 23, 116, 104, 105, 110, 103, 115, 32, 97, 114, 101, 32, 103, 111, 110, 110, 97, 32, 98, 101, 32, 98, 97, 100}
	if diff := cmp.Diff(Serialize(e), expected); diff != "" {
		t.Fatal(Serialize(e))
		t.Error("Serialize(e); (-got +want)", diff)
	}
}

func TestUnserialize(t *testing.T) {
	e := new(Error)
	e.Unserialize([]byte{6, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 57, 0, 4, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 56, 24, 103, 101, 116, 69, 114, 114, 111, 114, 70, 111, 114, 83, 101, 114, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110, 2, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 55, 23, 116, 104, 105, 110, 103, 115, 32, 97, 114, 101, 32, 103, 111, 110, 110, 97, 32, 98, 101, 32, 98, 97, 100})

	expected := getErrorForSerialization()
	if diff := cmp.Diff(e, expected); diff != "" {
		t.Error("Unserialize(); (-got +want)", diff)
	}
}

func TestIsKind(t *testing.T) {
	e := getErrorForSerialization()

	if diff := cmp.Diff(Is(e, bad), false); diff != "" {
		t.Error("Is(); (-got +want)", diff)
	}
	if diff := cmp.Diff(Is(e, worst), true); diff != "" {
		t.Error("Is(); (-got +want)", diff)
	}
}

func TestHas(t *testing.T) {
	var (
		Fine  = Kind(2)
		Worst = New(Kind(5))
		Worse = New(Kind(4), Worst)
		Bad   = New(Kind(3), Worse)
		x     = Bad
	)
	// yes
	e, b := Has(x, Kind(5))
	if diff := cmp.Diff(e, Worst); diff != "" {
		t.Error("Has(); (-got +want)", diff)
	}
	if diff := cmp.Diff(b, true); diff != "" {
		t.Error("Has(); (-got +want)", diff)
	}
	// no
	e, b = Has(x, Fine)
	if diff := cmp.Diff(e, nil); diff != "" {
		t.Error("Has(); (-got +want)", diff)
	}
	if diff := cmp.Diff(b, false); diff != "" {
		t.Error("Has(); (-got +want)", diff)
	}
}

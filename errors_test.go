package errors

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getErrorForSerialization() error {
	var TestingError Kind = 24
	return New("third error", Kind(TestingError), Here(), New("second error", New("first error"), Here()))
}

func TestEncode(t *testing.T) {
	e := getErrorForSerialization()

	expected := "third error\03724\037errors_test.go:12\036second error\0370\037errors_test.go:12\036first error\0370\037\036"
	if diff := cmp.Diff(Encode(e), expected); diff != "" {
		t.Error("Encode(e); (-got +want)", diff)
	}
}

func TestString(t *testing.T) {
	e := getErrorForSerialization()

	expected := `third error = 24 @ errors_test.go:12
	second error @ errors_test.go:12
	first error`
	if diff := cmp.Diff(e.Error(), expected); diff != "" {
		t.Fatal(e.Error())
		t.Error("Error.Error(); (-got +want)", diff)
	}
}

func TestJSONEncode(t *testing.T) {
	e := getErrorForSerialization()

	real, _ := json.Marshal(e)
	expected := []byte(`{"error":"third error","kind":24,"location":"errors_test.go:12","previous":{"error":"second error","kind":0,"location":"errors_test.go:12","previous":{"error":"first error","kind":0,"location":"","previous":null}}}`)
	if diff := cmp.Diff(real, expected); diff != "" {
		t.Error("json.Marshal(); (-got +want)", diff)
	}
}

func TestSerialize(t *testing.T) {
	e := getErrorForSerialization()

	expected := []byte{11, 116, 104, 105, 114, 100, 32, 101, 114, 114, 111, 114, 1, 24, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 50, 12, 115, 101, 99, 111, 110, 100, 32, 101, 114, 114, 111, 114, 1, 0, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 50, 11, 102, 105, 114, 115, 116, 32, 101, 114, 114, 111, 114, 1, 0, 0}
	if diff := cmp.Diff(Serialize(e), expected); diff != "" {
		t.Fatal(Serialize(e))
		t.Error("Serialize(e); (-got +want)", diff)
	}
}

func TestUnserialize(t *testing.T) {
	e := new(Error)
	e.Unserialize([]byte{11, 116, 104, 105, 114, 100, 32, 101, 114, 114, 111, 114, 1, 24, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 50, 12, 115, 101, 99, 111, 110, 100, 32, 101, 114, 114, 111, 114, 1, 0, 17, 101, 114, 114, 111, 114, 115, 95, 116, 101, 115, 116, 46, 103, 111, 58, 49, 50, 11, 102, 105, 114, 115, 116, 32, 101, 114, 114, 111, 114, 1, 0, 0})

	expected := getErrorForSerialization()
	if diff := cmp.Diff(e, expected); diff != "" {
		t.Error("Unserialize(); (-got +want)", diff)
	}
}

func TestIsKind(t *testing.T) {
	var TestingError Kind = 24
	e := getErrorForSerialization()

	expected := true
	if diff := cmp.Diff(IsKind(e, TestingError), expected); diff != "" {
		t.Error("IsKind(); (-got +want)", diff)
	}
}

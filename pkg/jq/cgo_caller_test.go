package jq

import (
	"testing"
)

func Test_CgoCall(t *testing.T) {

	in := `{"foo":"baz","bar":"quux"}`

	res, err := NewJq().Program(".").Run(in)
	if err != nil {
		t.Fatalf("expect Run not fail: %s", err)
	}
	if res != in {
		t.Fatalf("expect '%s', got '%s'", in, res)
	}

	if cgoCallsCh == nil {
		t.Fatalf("expect cgo calls channel should not be nil after first run")
	}

	res, err = NewJq().Program(".").Run(in)
	if err != nil {
		t.Fatalf("expect Run not fail: %s", err)
	}
	if res != in {
		t.Fatalf("expect '%s', got '%s'", in, res)
	}

	res, err = NewJq().Program(".").Run(in)
	if err != nil {
		t.Fatalf("expect Run not fail: %s", err)
	}
	if res != in {
		t.Fatalf("expect '%s', got '%s'", in, res)
	}
}

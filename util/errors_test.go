package util

import (
	"errors"
	"testing"
)

func TestErrors(t *testing.T) {

	err := new(Errors)
	err.Append(errors.New("Foo"))
	err.Append(errors.New("Bar"))

	if err.Error() != "Foo\nBar" {
		t.Fatalf("Expected 'Foo\nBar' got '%#v'", err.Error())
	}

}

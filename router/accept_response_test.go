package router

import "testing"

func TestAcceptResponse_String(t *testing.T) {
	response := &AcceptResponse{1.99, map[string]string{"indent": "true", "encoding": "utf,8"}, "application", "json"}
	expected := "application/json;q=1.99;encoding=\"utf,8\";indent=true"

	if response.String() != expected {
		t.Fatalf("expected '%s' got: '%s'", expected, response.String())
	}
}

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse_WriteHeader(t *testing.T) {
	response := newWrappedResponse(httptest.NewRecorder())
	response.WriteHeader(http.StatusBadGateway)

	if code := response.(*wrappedResponse).statusCode; code != http.StatusBadGateway {
		t.Fatalf("expected %d got %d", http.StatusBadGateway, code)
	} else {
		if response.(*wrappedResponse).statusLine() != "502 Bad Gateway" {
			t.Fatalf("expectexd '502 Bad Gateway' got '%s'", response.(*wrappedResponse).statusLine())
		}
	}
}

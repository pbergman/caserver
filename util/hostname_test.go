package util

import "testing"

func TestHostname_wildcard(t *testing.T) {
	if !ValidHost("foo.example.com", "*.example.com") {
		t.Fatal("Expected '*.example.com' to match 'foo.example.com'")
	}
	if !ValidHost("example.com", "*.example.com") {
		t.Fatal("Expected '*.example.com' to match 'example.com'")
	}
}

func TestHostname_name(t *testing.T) {
	if ValidHost("foo.example.com", "example.com") {
		t.Fatal("Not expected 'example.com' to match 'foo.example.com'")
	}

	if !ValidHost("example.com", "example.com") {
		t.Fatal("Expected 'example.com' to match 'example.com'")
	}
}

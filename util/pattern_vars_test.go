package util

import (
	"regexp"
	"testing"
)

func TestGetPatternVars(t *testing.T) {
	pattern := regexp.MustCompile(`^/api/v1/.+(?:\.(?P<extension>json|tar(?:\.gz)?|pem|text))$`)
	extensions := []string{"json", "tar", "tar.gz", "pem", "text"}
	for _, ext := range extensions {
		result := GetPatternVars("/api/v1/cert/aaaa."+ext, pattern)
		if result["extension"] != ext {
			t.Fatalf("expected to have %s got %s", ext, result["extension"])
		}
	}
}

func TestGetPatternVar(t *testing.T) {
	pattern := regexp.MustCompile(`^/api/v1/.+(?:\.(?P<extension>json|tar(?:\.gz)?|pem|text))$`)
	extensions := []string{"json", "tar", "tar.gz", "pem", "text"}
	for _, ext := range extensions {
		result := GetPatternVar("extension", "/api/v1/cert/aaaa."+ext, pattern)
		if result != ext {
			t.Fatalf("expected to have %s got %s", ext, result)
		}
	}
}

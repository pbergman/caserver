package util

import "testing"

func TestLowerSlug(t *testing.T) {

	list := map[string]string{
		"The quick brown fox jumps over the lazy dog":     "The_quick_brown_fox_jumps_over_the_lazy_dog",
		"The quick brown fox jumps over the lazy 10 dogs": "The_quick_brown_fox_jumps_over_the_lazy_10_dogs",
		"*.example.com":                                   "example_com",
		"*.*.*.":                                          "Ki4qLiou",
	}

	for v, k := range list {
		if ret := Slug(v); ret != k {
			t.Fatalf("expected '%s' got '%s'", k, ret)
		}
	}

}

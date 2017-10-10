package util

import "regexp"

var (
	// Domain names may be formed from the set of alphanumeric ASCII characters (a-z, A-Z, 0-9),
	// but characters are case-insensitive. In addition the hyphen is permitted if it is
	// surrounded by characters, digits or hyphens, although it is not to start or end a label.
	wildcard string = `(?:[a-z|0-9]+[a-z|0-9|\-]+[a-z|0-9]+\.)?`
)

func ValidHost(host, subject string) bool {
	var reg *regexp.Regexp

	if subject[0] == '*' && subject[1] == '.' {
		reg = regexp.MustCompile(`(?i)^` + wildcard + regexp.QuoteMeta(subject[2:]) + `?`)
	} else {
		// case-insensitive
		reg = regexp.MustCompile(`(?i)^` + regexp.QuoteMeta(subject) + `$`)
	}

	return reg.MatchString(host)
}

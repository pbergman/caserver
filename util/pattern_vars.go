package util

import "regexp"

func GetPatternVars(subject string, pattern *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	matches := pattern.FindStringSubmatch(subject)
	names := pattern.SubexpNames()
	for i, c := 1, len(names); i < c; i++ {
		result[names[i]] = matches[i]
	}
	return result
}

func GetPatternVar(key, subject string, pattern *regexp.Regexp) string {
	result := GetPatternVars(subject, pattern)
	if v, o := result[key]; o {
		return v
	}
	return ""
}

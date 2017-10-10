package util

import "encoding/base64"

func Slug(line string) string {
	ret := make([]rune, 0)
	for _, s := range line {
		switch {
		case s >= 'A' && s <= 'Z' || s >= 'a' && s <= 'z' || s >= '0' && s <= '9' || s == '-':
			ret = append(ret, s)
		default:
			if len(ret) > 0 {
				ret = append(ret, '_')
			}
		}
	}
	if len(ret) == 0 && len(line) > 0 {
		return base64.RawURLEncoding.EncodeToString([]byte(line))
	} else {
		return string(ret)
	}
}

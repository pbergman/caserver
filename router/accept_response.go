package router

import (
	"sort"
	"strconv"
	"strings"
)

type AcceptResponse struct {
	QualityValue float64
	Options      map[string]string
	Type         string
	SubType      string
}

func (c AcceptResponse) String() string {
	str := c.Type + "/" + c.SubType
	if c.QualityValue > 0 {
		str += ";q=" + strconv.FormatFloat(c.QualityValue, 'g', -1, 64)
	}
	keys := make([]string, 0)
	for key := range c.Options {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := c.Options[key]
		if strings.Contains(value, "\"") || strings.Contains(value, ",") || strings.Contains(value, ";") {
			str += ";" + key + "=\"" + strings.Replace(value, "\"", "\\\"", -1) + "\""
		} else {
			str += ";" + key + "=" + value
		}
	}
	return str
}

func (c AcceptResponse) GetType() string {
	return c.Type + "/" + c.SubType
}

func (c AcceptResponse) match(b *AcceptResponse) bool {
	return (c.Type == "*" && c.SubType == "*") || (b.Type == "*" && b.SubType == "*") || (c.Type == b.Type && (b.SubType == c.SubType || c.SubType == "*"))
}

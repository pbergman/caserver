package router

import (
	"sort"
	"strconv"
	"strings"
)

type AcceptResponses []*AcceptResponse

// implement of the sort.Interface
func (c AcceptResponses) Len() int           { return len(c) }
func (c AcceptResponses) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c AcceptResponses) Less(i, j int) bool { return c[i].QualityValue > c[j].QualityValue }

// Sort will sort the stack based on quality value (weight) of values
func (c *AcceptResponses) sort() { sort.Sort(c) }

// MatchFor will return 1 of the types that match based weight
// (quality-value) of registered types, when no content types
// is matched it will return 0
func (c *AcceptResponses) MatchFor(m ContentType) ContentType {
	c.sort()
	z := NewAcceptResponses(m.String())
	for d, e := 0, c.Len(); d < e; d++ {
		for a, b := 0, len(*z); a < b; a++ {
			if (*c)[d].match((*z)[a]) {
				return ContentTypeFromString((*z)[a].GetType())
			}
		}
	}
	return 0
}

// AcceptResponses will parse the given string as described in rfc2616 (section 14) given in an
// request as the Accept header to notify which content type client expect to return (or support)
//
// so for example:  text/html,application/xml;q=0.9;charset=utf-8,*/*;q=0.8
//
// will be parsed as:
//
// AcceptResponses{
//  []*AcceptResponse{
//     *AcceptResponse{
//          QualityValue:   1.0
//          Options:        {}
//          Type:           text
//          SubType:        html
//     },
//     *ContentType{
//          QualityValue:   0.9
//          Options:        {"charset": "utf-8"}
//          Type:           application
//          SubType:        xml
//     },
//     *ContentType{
//          QualityValue:   0.8
//          Options:        {}
//          Type:           *
//          SubType:        *
//     },
//  }
// }
func NewAcceptResponses(value string) *AcceptResponses {
	list := AcceptResponses(make([]*AcceptResponse, 0))
	for _, part := range strings.Split(value, ",") {
		parts := strings.Split(strings.TrimLeft(part, " "), ";")
		mime := strings.Split(parts[0], "/")
		item := &AcceptResponse{1, make(map[string]string), mime[0], mime[1]}
		for i, c := 1, len(parts); i < c; i++ {
			option := strings.Split(strings.TrimLeft(parts[i], " "), "=")
			if len(option) != 2 {
				continue
			}
			switch option[0] {
			case "q":
				// ignore errors
				if qv, o := strconv.ParseFloat(option[1], 64); o == nil {
					item.QualityValue = qv
				}
			default:
				first, last := option[1][0], option[1][len(option[1])-1]
				// check quoted, double quoted and escaped characters in string
				if first == last && (first == '"' || first == '\'') {
					option[1] = strings.Replace(option[1][1:len(option[1])-1], "\\"+string(first), string(first), -1)
				}
				item.Options[option[0]] = option[1]
			}
		}
		list = append(list, item)
	}
	return &list
}

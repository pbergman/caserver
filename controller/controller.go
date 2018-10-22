package controller

import (
	"regexp"

	"github.com/pbergman/caserver/router"
)

func newController(pattern string) Controller {
	return Controller{regexp.MustCompile(pattern)}
}

type Controller struct {
	pattern *regexp.Regexp
}

func (c Controller) Match(request *router.Request) bool {
	return c.pattern.MatchString(request.URL.Path)
}

func (c Controller) GetPathVars(request *router.Request) map[string]string {
	vars := make(map[string]string)
	matches := c.pattern.FindStringSubmatch(request.URL.Path)
	names := c.pattern.SubexpNames()
	for i, t := 1, len(names); i < t; i++ {
		vars[names[i]] = matches[i]
	}
	return vars
}

func (c Controller) GetPathVar(key string, request *router.Request) string {
	result := c.GetPathVars(request)
	if v, o := result[key]; o {
		return v
	}
	return ""
}

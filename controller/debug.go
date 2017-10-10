package controller

import (
	"net/http"
	"regexp"

	"github.com/pbergman/logger"
	"net/http/pprof"
)

type DebugController struct {
	pattern *regexp.Regexp
}

func (s DebugController) Handle(resp http.ResponseWriter, req *http.Request, logger logger.LoggerInterface) {
	match := s.pattern.FindStringSubmatch(req.RequestURI)

	switch match[1] {
	case "cmdline":
		pprof.Cmdline(resp, req)
		return
	case "profile":
		pprof.Profile(resp, req)
		return
	case "symbol":
		pprof.Symbol(resp, req)
		return
	case "trace":
		pprof.Trace(resp, req)
		return
	default:
		pprof.Index(resp, req)
		return
	}

}

func (s DebugController) Name() string {
	return "controller.debug"
}

func (s DebugController) Match(request *http.Request) bool {
	return s.pattern.MatchString(request.RequestURI)
}

func NewDebug() *DebugController {
	return &DebugController{pattern: regexp.MustCompile(`^/debug/pprof/([^$]+)`)}
}

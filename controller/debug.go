package controller

import (
	"net/http"

	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
	"net/http/pprof"
)

type DebugController struct {
	Controller
}

func (s DebugController) Handle(req *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {
	match := s.GetPathVar("type", req)
	switch match {
	case "cmdline":
		pprof.Cmdline(resp, req.Request)
		return
	case "profile":
		pprof.Profile(resp, req.Request)
		return
	case "symbol":
		pprof.Symbol(resp, req.Request)
		return
	case "trace":
		pprof.Trace(resp, req.Request)
		return
	default:
		pprof.Index(resp, req.Request)
		return
	}
}

func (s DebugController) Name() string {
	return "controller.debug"
}

func NewDebug() *DebugController {
	return &DebugController{newController(`^/debug/pprof/(?P<type>[^$]+)`)}
}

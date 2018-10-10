package controller

import (
	"net/http"
	"regexp"

	"github.com/pbergman/caserver/router"
	"github.com/pbergman/caserver/util"
	"github.com/pbergman/logger"
)

type PreAcceptHeader struct {
	pattern *regexp.Regexp
}

func (p PreAcceptHeader) Name() string {
	return "pre.hook.accept.header"
}

func (p PreAcceptHeader) Match(request *http.Request) bool {
	return p.pattern.MatchString(request.URL.Path)
}

func (p PreAcceptHeader) Handle(request *http.Request, logger logger.LoggerInterface) {
	extension := util.GetPatternVar("ext", request.URL.Path, p.pattern)

	request.URL.Path = request.URL.Path[:len(request.URL.Path)-len(extension)-1]
	logger.Debug("updating url path to: " + request.URL.Path)

	switch extension {
	case "json":
		request.Header.Set("accept", "application/json")
	case "tar":
		request.Header.Set("accept", "application/tar")
	case "tar.gz":
		request.Header.Set("accept", "application/tar+gzip")
	case "pem":
		request.Header.Set("accept", "application/pkix-cert")
	case "text", "txt":
		request.Header.Set("accept", "text/plain")
	}

	logger.Debug("set accept header to: " + request.Header.Get("accept"))
}

func NewPreAcceptHeaderHook() router.PreControllerInterface {
	return &PreAcceptHeader{
		pattern: regexp.MustCompile(`^/api/v1/.+(?:\.(?P<ext>json|tar(?:\.gz)?|pem|t(?:e)?xt))$`),
	}
}

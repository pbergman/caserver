package controller

import (
	"net/http"

	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

type PreAcceptHeader struct {
	Controller
}

func (p PreAcceptHeader) Name() string {
	return "pre.hook.accept.header"
}

func (p PreAcceptHeader) Handle(request *router.Request, header http.Header, logger logger.LoggerInterface) {
	extension := p.GetPathVar("ext", request)
	request.URL.Path = request.URL.Path[:len(request.URL.Path)-len(extension)-1]
	logger.Debug("updating url path to: " + request.URL.Path)
	switch extension {
	case "json":
		p.prefixAcceptHeader(request.Header, "application/json")
	case "tar":
		p.prefixAcceptHeader(request.Header, "application/tar")
	case "tar.gz":
		p.prefixAcceptHeader(request.Header, "application/tar+gzip")
	case "pem":
		p.prefixAcceptHeader(request.Header, "application/pkix-cert")
	case "text", "txt":
		p.prefixAcceptHeader(request.Header, "text/plain")
	}
	logger.Debug("set accept header to: '" + request.Header.Get("accept") + "'")
}

func (p PreAcceptHeader) prefixAcceptHeader(header http.Header, value string) {
	if accept := header.Get("accept"); accept != "" {
		value = value + ";q=9.0, " + accept // prefix string with an high weight
	}
	header.Set("accept", value)
}

func NewPreAcceptHeaderHook() router.PreControllerInterface {
	return &PreAcceptHeader{
		Controller: newController(`^/api/v1/.+(?:\.(?P<ext>json|tar(?:\.gz)?|pem|t(?:e)?xt))$`),
	}
}

package controller

import (
	"net/http"

	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

type PreResponseHeaders struct {
}

func (p PreResponseHeaders) Name() string {
	return "pre.hook.access.control.header"
}

func (p PreResponseHeaders) Match(request *router.Request) bool {
	return request.GetAcceptResponseType().MatchFor(router.ContentTypeAll) != router.ContentType(0)
}

func (p PreResponseHeaders) Handle(request *router.Request, header http.Header, logger logger.LoggerInterface) {
	switch request.GetAcceptResponseType().MatchFor(router.ContentTypeAll) {
	case router.ContentTypeJson:
		//logger.Debug("adding 'Access-Control-Allow-Origin: *' to the response")
		header.Add("access-Control-Allow-Origin", "*")
		header.Set("Content-Type", "application/json")
		header.Set("X-Content-Type-Options", "nosniff")
	case router.ContentTypeTar:
		header.Set("Content-Type", "application/tar")
		header.Set("X-Content-Type-Options", "nosniff")
	case router.ContentTypeTarGzip:
		header.Set("Content-Type", "application/tar+gzip")
		header.Set("X-Content-Type-Options", "nosniff")
	case router.ContentTypePkixCert:
		header.Set("Content-Type", "application/pkix-cert")
		header.Set("X-Content-Type-Options", "nosniff")
	case router.ContentTypeText:
		header.Set("Content-Type", "text/plain; charset=utf-8")
		header.Set("X-Content-Type-Options", "nosniff")
	}
}

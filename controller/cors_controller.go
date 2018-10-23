package controller

import (
	"net/http"

	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

// simple controller for the CORS options preflight
// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
type CorsController struct {}

func (c CorsController) Name() string {
	return "cors.control.header"
}

func (c CorsController) Match(request *router.Request) bool {
	return request.Method == "OPTIONS"
}

func (p CorsController) Handle(request *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {
	resp.Header().Set("Access-Control-Allow-Origin", request.Header.Get("Origin"))
	resp.Header().Set("Access-Control-Allow-Methods", request.Header.Get("Access-Control-Request-Method"))
	resp.Header().Set("Content-Type", "text/plain")
	resp.WriteHeader(http.StatusOK)
}

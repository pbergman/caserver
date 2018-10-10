package router

import (
	"net/http"

	"github.com/pbergman/logger"
)

type PreControllerInterface interface {
	Handle(*http.Request, http.Header, logger.LoggerInterface)
	Match(*http.Request) bool
	Name() string
}

type ControllerInterface interface {
	Handle(http.ResponseWriter, *http.Request, logger.LoggerInterface)
	Match(*http.Request) bool
	Name() string
}

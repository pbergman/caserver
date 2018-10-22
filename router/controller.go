package router

import (
	"net/http"

	"github.com/pbergman/logger"
)

type PreControllerInterface interface {
	Handle(*Request, http.Header, logger.LoggerInterface)
	Match(*Request) bool
	Name() string
}

type ControllerInterface interface {
	Handle(*Request, http.ResponseWriter, logger.LoggerInterface)
	Match(*Request) bool
	Name() string
}

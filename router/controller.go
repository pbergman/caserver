package router

import (
	"net/http"

	"github.com/pbergman/logger"
)

type ControllerInterface interface {
	Handle(http.ResponseWriter, *http.Request, logger.LoggerInterface)
	Match(*http.Request) bool
	Name() string
}

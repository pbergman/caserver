package controller

import (
	"github.com/pbergman/logger"
	"net/http"
)

func write_error(w http.ResponseWriter, error string, code int, log logger.LoggerInterface) {
	log.Error(error)
	http.Error(w, error, code)
}

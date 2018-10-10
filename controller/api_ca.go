package controller

import (
	"fmt"
	"net/http"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/logger"
)

type ApiCaController struct {
	manager *ca.Manager
}

func (s ApiCaController) Handle(resp http.ResponseWriter, req *http.Request, logger logger.LoggerInterface) {

	record := s.manager.Get(s.manager.GetCa())

	if req.Method != "GET" {
		write_error(resp, fmt.Sprintf("Method %s is not supported.", req.Method), http.StatusMethodNotAllowed, logger)
		return
	}

	if record == nil {
		write_error(resp, "Failed to find CA.", http.StatusInternalServerError, logger)
		return
	}

	// remove from record so wo`t be printed.
	record.SetPrivateKey(nil)

	if err := WriteResponse(resp, req, nil, record); err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
	}
}

func (s ApiCaController) Name() string {
	return "controller.api.ca"
}

func (s ApiCaController) Match(request *http.Request) bool {
	return request.URL.Path == "/api/v1/ca" && request.Method == "GET"
}

func NewApiCa(manager *ca.Manager) *ApiCaController {
	return &ApiCaController{
		manager: manager,
	}
}

package controller

import (
	"net/http"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

type ApiCertDeleteController struct {
	ApiCertController
}

func (a ApiCertDeleteController) Name() string {
	return "controller.api.cert.delete"
}

func (a ApiCertDeleteController) Match(request *router.Request) bool {
	return a.Controller.Match(request) && request.Method == "DELETE"
}

func NewApiCertDelete(manager *ca.Manager) *ApiCertDeleteController {
	return &ApiCertDeleteController{newApiCertController(manager, `^(?i)/api/v1/cert/(?P<id>[a-f0-9]{40})$`)}
}

func (a ApiCertDeleteController) Handle(req *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {
	id := a.GetPathVar("id", req)
	if record := a.manager.Lookup(id); record == nil {
		write_error(resp, "No record found for '"+id+"' .", http.StatusNotFound, logger)
	} else {
		if err := a.manager.Remove(record.GetId()); err != nil {
			write_error(resp, err.Error(), http.StatusNotFound, logger)
		} else {
			resp.WriteHeader(http.StatusAccepted)
		}
	}
}

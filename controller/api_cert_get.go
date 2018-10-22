package controller

import (
	"net/http"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

type ApiCertGetController struct {
	ApiCertController
}

func (a ApiCertGetController) Name() string {
	return "controller.api.cert.get"
}

func (a ApiCertGetController) Match(request *router.Request) bool {
	return a.Controller.Match(request) && request.Method == "GET"
}

func NewApiCertGet(manager *ca.Manager) *ApiCertGetController {
	return &ApiCertGetController{newApiCertController(manager, `^(?i)/api/v1/cert/(?P<id>[a-f0-9]{4,})$`)}
}

func (a ApiCertGetController) Handle(req *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {
	id := a.GetPathVar("id", req)
	if entry := a.manager.Lookup(id); entry == nil {
		write_error(resp, "could not find any record by "+id, http.StatusNotFound, logger)
		return
	} else {
		if entry.GetCertificate() == nil {
			write_error(resp, "no certificate found for record "+id, http.StatusNotFound, logger)
		} else {
			if err := WriteResponse(req, resp, a.getCa(), entry); err != nil {
				write_error(resp, err.Error(), http.StatusInternalServerError, logger)
			}
		}
	}
}

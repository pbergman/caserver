package controller

import (
	"net/http"
	"regexp"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/util"
	"github.com/pbergman/logger"
)

type ApiCertGetController struct {
	ApiCertController
}

func (a ApiCertGetController) Name() string {
	return "controller.api.cert.get"
}

func (a ApiCertGetController) Match(request *http.Request) bool {
	return a.pattern.MatchString(request.URL.Path) && request.Method == "GET"
}

func NewApiCertGet(manager *ca.Manager) *ApiCertGetController {
	return &ApiCertGetController{
		ApiCertController{
			pattern: regexp.MustCompile(`^(?i)/api/v1/cert/(?P<id>[a-f0-9]{4,})$`),
			manager: manager,
		},
	}
}

func (a ApiCertGetController) Handle(resp http.ResponseWriter, req *http.Request, logger logger.LoggerInterface) {
	id := util.GetPatternVar("id", req.URL.Path, a.pattern)
	if entry := a.manager.Lookup(id); entry == nil {
		write_error(resp, "could not find any record by "+id, http.StatusNotFound, logger)
		return
	} else {
		if entry.GetCertificate() == nil {
			write_error(resp, "no certificate found for record "+id, http.StatusNotFound, logger)
		} else {
			if err := WriteResponse(resp, req, a.getCa(), entry); err != nil {
				write_error(resp, err.Error(), http.StatusInternalServerError, logger)
			}
		}
	}
}

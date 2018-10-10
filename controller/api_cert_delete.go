package controller

import (
	"net/http"
	"regexp"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/util"
	"github.com/pbergman/logger"
)

type ApiCertDeleteController struct {
	ApiCertController
}

func (a ApiCertDeleteController) Name() string {
	return "controller.api.cert.delete"
}

func (a ApiCertDeleteController) Match(request *http.Request) bool {
	return a.pattern.MatchString(request.URL.Path) && request.Method == "DELETE"
}

func NewApiCertDelete(manager *ca.Manager) *ApiCertDeleteController {
	return &ApiCertDeleteController{
		ApiCertController{
			pattern: regexp.MustCompile(`^(?i)/api/v1/cert/(?P<id>[a-f0-9]{40})$`),
			manager: manager,
		},
	}
}

func (a ApiCertDeleteController) Handle(resp http.ResponseWriter, req *http.Request, logger logger.LoggerInterface) {
	id := util.GetPatternVar("id", req.URL.Path, a.pattern)
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

package controller

import (
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/router"
	"github.com/pbergman/logger"
)

type ApiCertCreateController struct {
	ApiCertController
}

func (a ApiCertCreateController) Name() string {
	return "controller.api.cert.create"
}

func (a ApiCertCreateController) Match(request *router.Request) bool {
	return a.Controller.Match(request) && request.Method == "POST"
}

func NewApiCertCreate(manager *ca.Manager) *ApiCertCreateController {
	return &ApiCertCreateController{newApiCertController(manager, `^(?i)/api/v1/cert$`)}
}

func (a ApiCertCreateController) Handle(req *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {

	record := a.getCa()

	if record == nil {
		write_error(resp, "Failed to find CA.", http.StatusInternalServerError, logger)
		return
	}

	if err := req.ParseForm(); err != nil {
		write_error(resp, err.Error(), http.StatusBadRequest, logger)
		return
	}

	subject, err := a.getSubject(req.Form)

	if err != nil {
		write_error(resp, err.Error(), http.StatusBadRequest, logger)
		return
	}

	if r := a.manager.Search(subject.CommonName); r != nil {
		resp.Header().Set("link", fmt.Sprintf("href=\"/api/v1/cert/%s\", rel=\"record\"", r.GetId().String()))
		write_error(resp, fmt.Sprintf("a csr exists for %s", subject.CommonName), http.StatusBadRequest, logger)
		return
	}

	var hosts []string

	if value, ok := req.Form["host"]; ok {
		hosts = value
	} else {
		hosts = []string{subject.CommonName}
	}

	entry, err := a.manager.NewCertificateRequest(hosts, subject, a.getBits(req))

	if err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
		return
	}

	if err := a.manager.SignCertificateRequest(entry, record); err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
		return
	}

	if err := WriteResponse(req, resp, record, entry); err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
		return
	}
}

func (a ApiCertCreateController) getBits(req *router.Request) int {
	var bits int = 2048

	if val, ok := req.Form["bits"]; ok {
		if v, err := strconv.Atoi(val[0]); err == nil {
			bits = v
		}
	}

	return bits
}

func (a ApiCertCreateController) getSubject(v url.Values) (name pkix.Name, err error) {
	for key, value := range v {
		switch strings.ToLower(key) {
		case "c", "country":
			name.Country = value
		case "o", "organization":
			name.Organization = value
		case "ou", "organizational_unit":
			name.OrganizationalUnit = value
		case "l", "locality":
			name.Locality = value
		case "p", "province":
			name.Province = value
		case "street", "street_address":
			name.StreetAddress = value
		case "postalcode", "postal_code":
			name.PostalCode = value
		case "cn", "common_name":
			name.CommonName = value[0]
		}
	}
	if name.CommonName == "" {
		err = errors.New("missing required 'cn' field")
	}
	return
}

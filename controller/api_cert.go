package controller

import (
	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/storage"
)

type ApiCertController struct {
	Controller
	manager *ca.Manager
}

func (a ApiCertController) getCa() storage.Record {
	return a.manager.Get(a.manager.GetCa())
}

func newApiCertController(manager *ca.Manager, pattern string) ApiCertController {
	return ApiCertController{
		manager:    manager,
		Controller: newController(pattern),
	}
}

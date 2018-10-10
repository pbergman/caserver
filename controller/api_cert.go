package controller

import (
	"archive/tar"
	"regexp"
	"time"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/storage"
)

type ApiCertController struct {
	pattern *regexp.Regexp
	manager *ca.Manager
}

func (a ApiCertController) getCa() storage.Record {
	return a.manager.Get(a.manager.GetCa())
}

func (s ApiCertController) tarHeader(name string, size int64) *tar.Header {
	return &tar.Header{
		Name:    name,
		Mode:    0600,
		Size:    size,
		ModTime: time.Now(),
	}
}

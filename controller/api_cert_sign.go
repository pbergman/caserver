package controller

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"net/http"

	"github.com/pbergman/caserver/ca"
	"github.com/pbergman/caserver/router"
	"github.com/pbergman/caserver/storage"
	"github.com/pbergman/logger"
)

type ApiCertSignController struct {
	ApiCertController
}

func (a ApiCertSignController) Name() string {
	return "controller.api.cert.sign"
}

func (a ApiCertSignController) Match(request *router.Request) bool {
	return a.Controller.Match(request) && request.Method == "PUT"
}

func NewApiCertSign(manager *ca.Manager) *ApiCertSignController {
	return &ApiCertSignController{newApiCertController(manager, `^(?i)/api/v1/cert$`)}
}

func (a ApiCertSignController) Handle(req *router.Request, resp http.ResponseWriter, logger logger.LoggerInterface) {
	file, _, err := req.FormFile("csr")

	if err != nil {
		if err == http.ErrMissingFile {
			write_error(resp, "missing required 'csr' post field", http.StatusBadRequest, logger)
			return
		} else {
			write_error(resp, err.Error(), http.StatusBadRequest, logger)
			return
		}
	}

	blockCsr, err := a.read(file)

	if err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
	}

	if blockCsr == nil {
		write_error(resp, "uploaded file was not a PEM encoded block.", http.StatusBadRequest, logger)
	}

	if blockCsr.Type != storage.BLOCK_TYPE_CSR {
		write_error(resp, "invalid PEM type", http.StatusBadRequest, logger)
	}

	caRecord := a.getCa()
	csr, err := x509.ParseCertificateRequest(blockCsr.Bytes)

	if err != nil {
		write_error(resp, err.Error(), http.StatusBadRequest, logger)
	}

	if csr.Subject.CommonName == "" {
		write_error(resp, "missing required 'cn' field in csr", http.StatusBadRequest, logger)
	}

	cer, err := a.manager.GetFactory().NewCertificate(csr, caRecord.GetCertificate(), caRecord.GetPrivateKey())

	if err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
	}

	cerRecord := a.manager.NewRecord()
	cerRecord.SetCertificate(cer)

	if err := WriteResponse(req, resp, caRecord, cerRecord); err != nil {
		write_error(resp, err.Error(), http.StatusInternalServerError, logger)
	}
}

func (s ApiCertSignController) read(reader io.Reader) (*pem.Block, error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}
	raw, buf := make([]byte, 0), make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		raw = append(raw, buf[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}
	block, _ := pem.Decode(raw)
	return block, nil
}

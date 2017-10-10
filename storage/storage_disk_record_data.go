package storage

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
)

// DiskRecordHeader is the data part of the record (DiskRecord)
type DiskRecordData struct {
	key *rsa.PrivateKey
	pem *x509.Certificate
	csr *x509.CertificateRequest
}

func (d DiskRecordData) GetPrivateKey() *rsa.PrivateKey {
	return d.key
}

func (d DiskRecordData) GetCertificate() *x509.Certificate {
	return d.pem
}

func (d DiskRecordData) GetCertificateRequest() *x509.CertificateRequest {
	return d.csr
}

func (d DiskRecordData) HasPrivateKey() bool {
	return nil != d.key
}

func (d DiskRecordData) HasCertificate() bool {
	return nil != d.pem
}

func (d DiskRecordData) HasCertificateRequest() bool {
	return nil != d.csr
}

func (d DiskRecordData) WritePrivateKey(w io.Writer) error {
	return d.export(d.key, w)
}

func (d DiskRecordData) WriteCertificate(w io.Writer) error {
	return d.export(d.pem, w)
}

func (d DiskRecordData) WriteCertificateRequest(w io.Writer) error {
	return d.export(d.csr, w)
}

func (d DiskRecordData) export(v interface{}, w io.Writer) error {
	var block *pem.Block
	switch t := v.(type) {
	case *rsa.PrivateKey:
		block = &pem.Block{
			Type:  BLOCK_TYPE_KEY,
			Bytes: x509.MarshalPKCS1PrivateKey(t),
		}
	case *x509.Certificate:
		block = &pem.Block{
			Type:  BLOCK_TYPE_CER,
			Bytes: t.Raw,
		}
	case *x509.CertificateRequest:
		block = &pem.Block{
			Type:  BLOCK_TYPE_CSR,
			Bytes: t.Raw,
		}
	case nil:
		return fmt.Errorf("could not export a nil value")
	}
	return pem.Encode(w, block)
}

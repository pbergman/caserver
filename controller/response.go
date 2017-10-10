package controller

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pbergman/caserver/storage"
	"github.com/pbergman/caserver/util"
)

type ctid uint16

const (
	CONTENT_TYPE_TEXT ctid = 1 << iota
	CONTENT_TYPE_JSON
	CONTENT_TYPE_TAR
	CONTENT_TYPE_TAR_GZIP
	CONTENT_TYPE_PKIX_CERT

	CONTENT_TYPE_ALL ctid = CONTENT_TYPE_TEXT | CONTENT_TYPE_JSON | CONTENT_TYPE_TAR | CONTENT_TYPE_TAR_GZIP | CONTENT_TYPE_PKIX_CERT
)

// ReadAndSetContentType will read the accept header, set the appropriate
// response header and return the matching CONTENT_TYPE_* const.
func readAndSetContentType(resp http.ResponseWriter, req *http.Request, allowed ctid) ctid {
	accept := strings.Split(req.Header.Get("accept"), ",")
	for i, c := 0, len(accept); i < c; i++ {
		if index := strings.IndexByte(accept[i], ';'); index >= 0 {
			accept[i] = accept[i][:index]
		}
		switch {
		case (CONTENT_TYPE_JSON == (CONTENT_TYPE_JSON & allowed)) && accept[i] == "application/json":
			resp.Header().Set("Content-Type", "application/json")
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			return CONTENT_TYPE_JSON
		case (CONTENT_TYPE_TAR == (CONTENT_TYPE_TAR & allowed)) && accept[i] == "application/tar":
			resp.Header().Set("Content-Type", "application/tar")
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			return CONTENT_TYPE_TAR
		case (CONTENT_TYPE_TAR_GZIP == (CONTENT_TYPE_TAR_GZIP & allowed)) && accept[i] == "application/tar+gzip":
			resp.Header().Set("Content-Type", "application/tar+gzip")
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			return CONTENT_TYPE_TAR_GZIP
		case (CONTENT_TYPE_PKIX_CERT == (CONTENT_TYPE_PKIX_CERT & allowed)) && accept[i] == "application/pkix-cert":
			resp.Header().Set("Content-Type", "application/pkix-cert")
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			return CONTENT_TYPE_PKIX_CERT
		case (CONTENT_TYPE_TEXT == (CONTENT_TYPE_TEXT & allowed)) && (accept[i] == "text/plain" || accept[i] == "*/*"):
			resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
			resp.Header().Set("X-Content-Type-Options", "nosniff")
			return CONTENT_TYPE_TEXT
		}
	}

	return ctid(0)
}

func tarFileHeader(name string, size int64) *tar.Header {
	return &tar.Header{
		Name:    name,
		Mode:    0600,
		Size:    size,
		ModTime: time.Now(),
	}
}

// writeTar writes a tar file of available fields in record
func writeTarResponse(writer io.Writer, ca, cer storage.Record) error {
	if !cer.HasCertificate() {
		return errors.New("can not write a tar file without a certificate")
	}
	name := util.Slug(cer.GetCertificate().Subject.CommonName)
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()
	// add private key when available
	if cer.HasPrivateKey() {
		if err := tarWriter.WriteHeader(tarFileHeader(name+".key", cer.BlockKeyLen())); err != nil {
			return err
		}
		if err := cer.WritePrivateKey(tarWriter); err != nil {
			return err
		}
	}
	// add certificate request when available
	if cer.HasCertificateRequest() {
		if err := tarWriter.WriteHeader(tarFileHeader(name+".csr", cer.BlockCsrLen())); err != nil {
			return err
		}
		if err := cer.WriteCertificateRequest(tarWriter); err != nil {
			return err
		}
	}
	if cer.HasCertificate() {
		var size int64
		if ca != nil {
			size = cer.BlockPemLen() + ca.BlockPemLen()
		} else {
			size = cer.BlockPemLen()
		}
		// add certificate
		if err := tarWriter.WriteHeader(tarFileHeader(name+".pem", size)); err != nil {
			return err
		}
		if err := cer.WriteCertificate(tarWriter); err != nil {
			return err
		}
		// add chained ca certificate when available
		if ca != nil {
			if err := ca.WriteCertificate(tarWriter); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeTarGzResponse will add gzip compression to the tar writer
func writeTarGzResponse(writer io.Writer, ca, record storage.Record) error {
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()
	return writeTarResponse(gzipWriter, ca, record)
}

// writeTextResponse will key, pem and csr (if available) to writer
func writeTextResponse(writer io.Writer, ca, record storage.Record) error {
	if record.HasPrivateKey() {
		if err := record.WritePrivateKey(writer); err != nil {
			return nil
		}
	}
	if record.HasCertificateRequest() {
		if err := record.WriteCertificateRequest(writer); err != nil {
			return nil
		}
	}
	if record.HasCertificate() {
		if err := record.WriteCertificate(writer); err != nil {
			return nil
		}
		if ca != nil {
			if err := ca.WriteCertificate(writer); err != nil {
				return nil
			}
		}
	}
	return nil
}

// writeJsonResponse will write a json response
func writeJsonResponse(writer io.Writer, ca, record storage.Record) error {
	data := make(map[string]interface{}, 0)
	buf := new(bytes.Buffer)
	if record.HasPrivateKey() {
		if err := record.WritePrivateKey(buf); err != nil {
			return err
		}
		data["key"] = buf.String()
		buf.Reset()
	}
	if record.HasCertificateRequest() {
		if err := record.WriteCertificateRequest(buf); err != nil {
			return err
		}
		data["csr"] = buf.String()
		buf.Reset()
	}
	if record.HasCertificate() {
		if err := record.WriteCertificate(buf); err != nil {
			return err
		}
		if ca != nil {
			if err := ca.WriteCertificate(buf); err != nil {
				return err
			}
		}
		data["pem"] = buf.String()
		buf.Reset()
	}
	enc := json.NewEncoder(writer)
	enc.SetIndent("", " ")
	return enc.Encode(data)
}

func nameFromRecord(record storage.Record) string {
	if record.GetId() != nil {
		return record.GetId().String()
	}
	if csr := record.GetCertificateRequest(); csr != nil {
		if csr.Subject.CommonName != "" {
			return util.Slug(csr.Subject.CommonName)
		}
	}
	if pem := record.GetCertificate(); pem != nil {
		if pem.Subject.CommonName != "" {
			return util.Slug(pem.Subject.CommonName)
		}
	}
	buf, _ := json.Marshal(record)
	return hex.EncodeToString(buf)
}

func WriteResponse(resp http.ResponseWriter, req *http.Request, caRecord, cerRecord storage.Record) error {
	name := nameFromRecord(cerRecord)
	switch readAndSetContentType(resp, req, CONTENT_TYPE_ALL) {
	case CONTENT_TYPE_JSON:
		resp.Header().Set("Content-Disposition", "inline; filename=\""+name+".json\"")
		return writeJsonResponse(resp, caRecord, cerRecord)
	case CONTENT_TYPE_TEXT:
		return writeTextResponse(resp, caRecord, cerRecord)
	case CONTENT_TYPE_TAR_GZIP:
		resp.Header().Set("Content-Disposition", "inline; filename=\""+name+".tar.gz\"")
		return writeTarGzResponse(resp, caRecord, cerRecord)
	case CONTENT_TYPE_TAR:
		resp.Header().Set("Content-Disposition", "inline; filename=\""+name+".tar\"")
		return writeTarResponse(resp, caRecord, cerRecord)
	case CONTENT_TYPE_PKIX_CERT:
		resp.Header().Set("Content-Disposition", "inline; filename=\""+name+".pem\"")
		return writeTextResponse(resp, caRecord, cerRecord)
	default:
		resp.WriteHeader(http.StatusNotAcceptable)
	}
	return nil
}

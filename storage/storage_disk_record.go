package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"io"

	"github.com/pbergman/caserver/util"
)

const (
	BLOCK_TYPE_KEY string = "RSA PRIVATE KEY"
	BLOCK_TYPE_CER string = "CERTIFICATE"
	BLOCK_TYPE_CSR string = "CERTIFICATE REQUEST"
)

func NewDiskRecord(s *DiskStorage, k *StorageKey) *DiskRecord {
	return &DiskRecord{
		DiskRecordData{},
		DiskRecordHeader{
			storage: s,
			id:      k,
		},
	}
}

// DiskRecord is instance of Storage and provide a key storage backed by files.
type DiskRecord struct {
	DiskRecordData
	DiskRecordHeader
}

// MarshalBinary will convert the struct to a custom binary stream and every record
// will be signed so that the UnmarshalBinary function can validate the data.
func (d DiskRecord) MarshalBinary() (data []byte, err error) {

	if d.storage == nil {
		return nil, errors.New("could not marshal data without a storage reference")
	}

	buf := new(bytes.Buffer)
	mac := hmac.New(sha256.New, d.storage.key[:])
	writer := io.MultiWriter(mac, buf)

	_, err = writer.Write([]byte{
		byte(d.mode),
		byte(d.size_key),
		byte(d.size_key >> 8),
		byte(d.size_pem),
		byte(d.size_pem >> 8),
		byte(d.size_csr),
		byte(d.size_csr >> 8),
	})

	if err != nil {
		return nil, err
	}

	if d.size_key > 0 {
		if _, err := writer.Write(x509.MarshalPKCS1PrivateKey(d.key)); err != nil {
			return nil, err
		}
	}

	if d.size_pem > 0 {
		if _, err := writer.Write(d.pem.Raw); err != nil {
			return nil, err
		}
	}

	if d.size_csr > 0 {
		if _, err := writer.Write(d.csr.Raw); err != nil {
			return nil, err
		}
	}

	return append(mac.Sum(nil), buf.Bytes()...), nil
}

// UnmarshalBinary a custom implementation for the gob.Decoder, it will
// also validate the signature and return a error if that fails.
func (d *DiskRecord) UnmarshalBinary(data []byte) error {

	sig, data := data[:sha256.Size], data[sha256.Size:]
	mac := hmac.New(sha256.New, d.storage.key[:])
	mac.Write(data)

	// validate
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return errors.New("invalid record")
	}

	head, data := data[:7], data[7:]

	d.mode = head[0]
	d.size_key = int(head[1]) | int(head[2])<<8
	d.size_pem = int(head[3]) | int(head[4])<<8
	d.size_csr = int(head[5]) | int(head[6])<<8

	var raw []byte
	var err error

	if d.size_key > 0 {
		raw, data = data[:d.size_key], data[d.size_key:]
		if d.key, err = x509.ParsePKCS1PrivateKey(raw); err != nil {
			return err
		}
	}

	if d.size_pem > 0 {
		raw, data = data[:d.size_pem], data[d.size_pem:]
		if d.pem, err = x509.ParseCertificate(raw); err != nil {
			return err
		}
	}

	if d.size_csr > 0 {
		raw, data = data[:d.size_csr], data[d.size_csr:]
		if d.csr, err = x509.ParseCertificateRequest(raw); err != nil {
			return err
		}
	}

	return nil
}

func (d *DiskRecord) SetPrivateKey(key *rsa.PrivateKey) {
	if key != nil {
		d.key = key
		d.size_key = len(x509.MarshalPKCS1PrivateKey(key))
	} else {
		d.key = nil
		d.size_key = 0
	}
}

func (d *DiskRecord) SetCertificate(cert *x509.Certificate) {
	if cert != nil {
		d.pem = cert
		d.size_pem = len(cert.Raw)
	} else {
		d.pem = nil
		d.size_pem = 0
	}
}

func (d *DiskRecord) SetCertificateRequest(cert *x509.CertificateRequest) {
	if cert != nil {
		d.csr = cert
		d.size_csr = len(cert.Raw)
	} else {
		d.csr = nil
		d.size_csr = 0
	}
}

// BlockPemLen will calculate the size of a generated pem block.
func (d DiskRecord) BlockPemLen() int64 {
	return util.PemLength(d.size_pem, BLOCK_TYPE_CER)
}

func (d DiskRecord) BlockCsrLen() int64 {
	return util.PemLength(d.size_csr, BLOCK_TYPE_CSR)
}

func (d DiskRecord) BlockKeyLen() int64 {
	return util.PemLength(d.size_key, BLOCK_TYPE_KEY)
}

func (d DiskRecord) IsCa() bool {
	return d.isCa()
}

func (d DiskRecord) GetId() *StorageKey {
	return d.id
}

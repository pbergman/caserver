package storage

import (
	"crypto/rsa"
	"crypto/x509"
	"io"
)

type Record interface {
	GetId() *StorageKey
	IsCa() bool
	// getter
	GetPrivateKey() *rsa.PrivateKey
	GetCertificate() *x509.Certificate
	GetCertificateRequest() *x509.CertificateRequest
	// setters
	SetPrivateKey(*rsa.PrivateKey)
	SetCertificate(*x509.Certificate)
	SetCertificateRequest(*x509.CertificateRequest)
	// exporters
	WritePrivateKey(io.Writer) error
	WriteCertificate(io.Writer) error
	WriteCertificateRequest(io.Writer) error
	// calculations for compression archives
	BlockPemLen() int64
	BlockCsrLen() int64
	BlockKeyLen() int64
	// checks
	HasPrivateKey() bool
	HasCertificate() bool
	HasCertificateRequest() bool
}

type Storage interface {
	// Persist will save the record to storage
	Persist(record Record) (*StorageKey, error)
	// Open will search a record by given kid
	// if not found will return nil
	Open(*StorageKey) (Record, error)
	// Lookup will do lookup for a short id (storage key)
	// when multiple matches it will return the first.
	Lookup(id string) (Record, error)
	// Search will search for every record and check if the
	// set CN will match. If no match nil will be returned.
	Search(string) Record
	// Remove a entry based on given kid
	Remove(*StorageKey) error
	// Has will check if record exist for given kid
	Has(*StorageKey) bool
	// GetCa will return list of keys that represent a ca
	GetCa() []*StorageKey
	// will return a new record based on storage type
	NewRecord() Record
	// Each will walk trough all records or til false is returned
	Each(call func(Record) bool) error
}

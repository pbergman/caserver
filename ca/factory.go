package ca

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net"
	"os"
	"time"
)

type factory struct {
	// pem (pem not after) represents 3
	// int`t for year, month and day that
	// will be added to current time and
	// set as NotAfter in the Certificate.
	pna *[3]int
	// same as pna but for the CA certificate
	cna *[3]int
	// the base serial used for creating new
	// serial numbers
	serial *big.Int
}

func (f *factory) init() {
	if f.serial == nil {
		f.serial = new(big.Int)
	}
	// create a base id from the machine id
	if file, err := os.Open("/etc/machine-id"); err == nil {
		defer file.Close()
		// should corresponds to a 16-byte value, that would be 32 byte hex string.
		// see https://www.freedesktop.org/software/systemd/man/machine-id.html
		buf := make([]byte, 32)
		n, _ := file.Read(buf)
		f.serial.SetBytes(buf[:n])
	}
}

// NewCertificateAuthority creates Certificate Authority using the
// given private key and returns a certificate in DER encoding.
func (f factory) NewCertificateAuthority(key *rsa.PrivateKey, subject pkix.Name) (*x509.Certificate, error) {
	f.checkSubject(&subject)
	ski, err := f.createSubjectKeyId(key.PublicKey)
	if err != nil {
		return nil, err
	}
	tmpl := x509.Certificate{
		SerialNumber:          f.newSerialNumber(nil),
		Subject:               subject,
		NotBefore:             time.Now().Add(-600).UTC(),
		NotAfter:              time.Now().AddDate((*f.cna)[0], (*f.cna)[1], (*f.cna)[2]).UTC(),
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:         true,
		MaxPathLen:   0,
		SubjectKeyId: ski,
	}
	raw, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(raw)
}

func (f factory) NewCertificateRequest(key *rsa.PrivateKey, subject pkix.Name, hosts []string) (*x509.CertificateRequest, error) {
	f.checkSubject(&subject)
	tmpl := &x509.CertificateRequest{Subject: subject}
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
		} else {
			tmpl.DNSNames = append(tmpl.DNSNames, host)
		}
	}
	raw, err := x509.CreateCertificateRequest(rand.Reader, tmpl, key)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificateRequest(raw)
}

func (f factory) NewCertificate(csr *x509.CertificateRequest, caCert *x509.Certificate, caKey *rsa.PrivateKey) (*x509.Certificate, error) {
	f.checkSubject(&csr.Subject)
	ski, err := f.createSubjectKeyId(*csr.PublicKey.(*rsa.PublicKey))
	if err != nil {
		return nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: f.newSerialNumber(caCert),
		Subject:      csr.Subject,
		NotBefore:    time.Now().Add(-600).UTC(),
		NotAfter:     time.Now().AddDate((*f.pna)[0], (*f.pna)[1], (*f.pna)[2]).UTC(),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SubjectKeyId: ski,
		DNSNames:     csr.DNSNames,
		IPAddresses:  csr.IPAddresses,
	}
	raw, err := x509.CreateCertificate(rand.Reader, tmpl, caCert, csr.PublicKey, caKey)
	if err != nil {
		return nil, err
	}
	return x509.ParseCertificate(raw)
}

func (a *factory) checkSubject(name *pkix.Name) {
	if name.SerialNumber == "" {
		if buf, err := json.Marshal(name); err == nil {
			hash := sha1.Sum(buf)
			name.SerialNumber = hex.EncodeToString(hash[:])
		}
	}
}

// newSerialNumber will generate a unique serial number for a certificate
// based on the machine id (when available), time from now and root pem
func (a *factory) newSerialNumber(root *x509.Certificate) *big.Int {
	serial := new(big.Int)
	serial.SetBytes(a.serial.Bytes())
	if root != nil {
		return serial.Add(big.NewInt((time.Now().Unix()+int64(1))-root.NotBefore.UnixNano()), serial)
	} else {
		return serial.Add(big.NewInt(time.Now().Unix()), serial)
	}
}

// createSubjectKeyId will create a byte slice that represents
// a (SHA-1 hash value) ASN.1 encoding of a public key.
func (f factory) createSubjectKeyId(key rsa.PublicKey) ([]byte, error) {
	buf, err := asn1.Marshal(key)
	if err != nil {
		return nil, err
	}
	hasher := sha1.New()
	hasher.Write(buf)
	return hasher.Sum(nil), nil
}

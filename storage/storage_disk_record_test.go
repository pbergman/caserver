package storage

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
)

func newTestCa(f *factory, t *testing.T) (*rsa.PrivateKey, *x509.Certificate) {
	key, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		t.Fatal(err)
	}

	cer, err := f.NewCertificateAuthority(key, pkix.Name{CommonName: "example CA"})

	if err != nil {
		t.Fatal(err)
	}

	return key, cer
}

func TestDiskRecord_PemLen(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		t.Fatal(err)
	}

	factory := new(factory)
	// new CA key and certificate
	ck, cc := newTestCa(factory, t)

	csr, err := factory.NewCertificateRequest(key, pkix.Name{CommonName: "example"}, []string{"example.com"})

	if err != nil {
		t.Fatal(err)
	}

	cer, err := factory.NewCertificate(csr, cc, ck)

	if err != nil {
		t.Fatal(err)
	}

	record := NewDiskRecord(new(DiskStorage), nil)
	record.SetCertificate(cer)
	record.SetCertificateRequest(csr)
	record.SetPrivateKey(key)

	buf := new(bytes.Buffer)

	buf.Reset()
	record.WriteCertificate(buf)

	if buf.Len() != int(record.BlockPemLen()) {
		t.Fatalf("Expected %d got %d", buf.Len(), record.BlockPemLen())
	}

	buf.Reset()
	record.WriteCertificateRequest(buf)

	if buf.Len() != int(record.BlockCsrLen()) {
		t.Fatalf("Expected %d got %d", buf.Len(), record.BlockCsrLen())
	}

	buf.Reset()
	record.WritePrivateKey(buf)

	if buf.Len() != int(record.BlockKeyLen()) {
		t.Fatalf("Expected %d got %d", buf.Len(), record.BlockKeyLen())
	}
}

func TestDiskRecord_Marshal(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		t.Fatal(err)
	}

	factory := new(factory)

	ck, cc := newTestCa(factory, t)

	csr, err := factory.NewCertificateRequest(key, pkix.Name{CommonName: "TEST", Country: []string{"NL"}}, []string{"example.com", "*.example.com"})

	if err != nil {
		t.Fatal(err)
	}

	cer, err := factory.NewCertificate(csr, cc, ck)

	if err != nil {
		t.Fatal(err)
	}

	record := NewDiskRecord(new(DiskStorage), nil)
	record.SetCertificate(cer)
	record.SetCertificateRequest(csr)
	record.SetPrivateKey(key)

	buf, err := record.MarshalBinary()

	if err != nil {
		t.Fatal(err)
	}

	newRecord := NewDiskRecord(new(DiskStorage), nil)

	if err := newRecord.UnmarshalBinary(buf); err != nil {
		t.Fatal(err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(cc)

	opts := x509.VerifyOptions{
		DNSName: "foo.example.com",
		Roots:   roots,
	}

	if _, err := newRecord.GetCertificate().Verify(opts); err != nil {
		t.Fatal(err)
	}
}

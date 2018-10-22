package ca

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/pbergman/caserver/storage"
)

/**
 * moved the disk storage test to the ca packages because of import
 * cycle and because we are testing only the public api of the disk
 * storage this is the best solution for now...
 */

func newTestCa(f FactoryInterface, t *testing.T) (*rsa.PrivateKey, *x509.Certificate) {
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

func newDiskStorage() *storage.DiskStorage {
	key := new([32]byte)
	rand.Read(key[:])
	return storage.NewDiskStorage("/tmp/", key)
}

func TestDiskRecord_PemLen(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)

	if err != nil {
		t.Fatal(err)
	}

	factory := NewFactory([3]int{1, 2, 3}, [3]int{4, 5, 6}, nil)
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

	record := storage.NewDiskRecord(newDiskStorage(), nil)
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

	factory := NewFactory([3]int{1, 2, 3}, [3]int{4, 5, 6}, nil)

	ck, cc := newTestCa(factory, t)

	csr, err := factory.NewCertificateRequest(key, pkix.Name{CommonName: "TEST", Country: []string{"NL"}}, []string{"example.com", "*.example.com"})

	if err != nil {
		t.Fatal(err)
	}

	cer, err := factory.NewCertificate(csr, cc, ck)

	if err != nil {
		t.Fatal(err)
	}

	ds := newDiskStorage()

	record := storage.NewDiskRecord(ds, nil)
	record.SetCertificate(cer)
	record.SetCertificateRequest(csr)
	record.SetPrivateKey(key)

	buf, err := record.MarshalBinary()

	if err != nil {
		t.Fatal(err)
	}

	newRecord := storage.NewDiskRecord(ds, nil)

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

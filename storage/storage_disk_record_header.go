package storage

const (
	MODE_IS_CA uint8 = (1 << iota)
)

// DiskRecordHeader is the header part of the record (DiskRecord)
type DiskRecordHeader struct {
	// a mode bit that hold information like it has a private
	// key, certificate, certificate request or is a CA.
	mode uint8
	// the message signature
	signature [20]byte
	// size bits of key, certificate and certificate
	// request (in order that the are marshaled)
	size_key int
	size_pem int
	size_csr int
	// a ref back to the storage manager used to get the key
	// for signing and verifying the signature.
	storage *DiskStorage
	// when file opened, the key will be added so when updating
	// it will know the old location. This because the key is
	// based on the content.
	id *StorageKey
}

func (h DiskRecordHeader) isCa() bool {
	return MODE_IS_CA == (MODE_IS_CA & h.mode)
}

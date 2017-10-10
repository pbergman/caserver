package storage

import "encoding/hex"

type StorageKey [20]byte

func (s StorageKey) String() string {
	return hex.EncodeToString(s[:])
}

func (s StorageKey) Bytes() []byte {
	return s[:]
}

func NewStorageKeyFromString(s string) *StorageKey {
	if buf, err := hex.DecodeString(s); err != nil {
		return nil
	} else {
		return NewStorageKeyFromBytes(buf)
	}
}

func NewStorageKeyFromBytes(b []byte) *StorageKey {
	key := new(StorageKey)
	copy(key[:], b[:20])
	return key
}

package storage

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestStorageKey(t *testing.T) {

	buf := []byte{
		0x71, 0x05, 0x7c, 0x93, 0x46, 0x11, 0x44, 0x0b, 0xdf, 0xf0,
		0xe9, 0x36, 0xfa, 0x52, 0x85, 0x2c, 0x21, 0xa6, 0x8c, 0xc2,
		0x08, 0x9b, 0xd2, 0xfa, 0xe0, 0x19, 0xa2, 0x77, 0x7e, 0x87,
		0xaf, 0x59, 0xe6, 0x3e, 0x59, 0x20, 0x17, 0x78, 0xcc, 0x73,
	}

	kid := NewStorageKeyFromBytes(buf)

	if !bytes.Equal(kid[:], buf[:20]) {
		t.Fatalf("Expected:\n%sGot:\n%s", hex.Dump(buf[:20]), hex.Dump(kid[:]))
	}

	if kid.String() != hex.EncodeToString(buf[:20]) {
		t.Fatalf("Expected: '%s' Got: '%s'", kid.String(), hex.EncodeToString(buf[:20]))
	}

	new_kid := NewStorageKeyFromString(hex.EncodeToString(buf[:20]))

	if !bytes.Equal(kid[:], new_kid[:]) {
		t.Fatalf("Expected:\n%sGot:\n%s", hex.Dump(kid[:]), hex.Dump(new_kid[:]))
	}
}

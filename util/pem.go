package util

import (
	"encoding/base64"
	"math"
)

// PemLength will calculate a pem block size
func PemLength(n int, s string) int64 {
	if n > 0 {
		size := base64.StdEncoding.EncodedLen(n)            // basic length of raw data
		size += int(math.Ceil(float64(size) / float64(64))) // every line has a newline
		size += 17 + len(s)                                 // `-----BEGIN BLOCK_TYPE-----\n`
		size += 15 + len(s)                                 // `-----END BLOCK_TYPE-----\n`
		return int64(size)
	} else {
		return 0
	}
}

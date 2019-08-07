package bs_test

import (
	"crypto/rand"
)

// Common test functions

func generateRandomData(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

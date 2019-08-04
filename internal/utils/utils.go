package utils

import (
	"io"
)

// Must suffice most cases
const maxMemoryForReader = 1 << 20 // 1MB

// GetReaderSize returns a copy of the io.Reader and the size of this reader (original io.Reader should contain small amount of data)
//
// It panics if io.Copy finished with an error.
func GetReaderSize(r io.Reader) (newReader io.Reader, size int64) {
	b := NewBuffer(maxMemoryForReader)
	size, err := io.Copy(b, r)
	if err != nil {
		panic(err)
	}
	b.Finish()

	return b, size
}

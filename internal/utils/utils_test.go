package utils

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestGetReaderSize(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		data []byte
		size int64
	}{
		{
			data: []byte("Hello!"),
			size: 6,
		},
		{
			data: []byte(""),
			size: 0,
		},
		{
			data: []byte("Привет, мир!"),
			size: 21, // 9 * 2 + 3
		},
	}

	for i, tt := range tests {
		b := &bytes.Buffer{}
		b.Write(tt.data)

		var r io.Reader = b
		var size int64
		r, size = GetReaderSize(r)

		assert.Equalf(tt.size, size, "Test #%d: different sizes", i+1)

		data, _ := ioutil.ReadAll(r)
		assert.Equalf(tt.data, data, "Test #%d: different data were read", i+1)
	}
}

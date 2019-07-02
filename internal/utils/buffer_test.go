package utils

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_CheckBufferAndFileSize(t *testing.T) {
	tests := []struct {
		maxSize int
		//
		data []byte
		//
		bufferSize int
		fileSize   int
	}{
		{
			maxSize:    15,
			data:       make([]byte, 10),
			bufferSize: 10,
			fileSize:   0,
		},
		{
			maxSize:    15,
			data:       make([]byte, 15),
			bufferSize: 15,
			fileSize:   0,
		},
		{
			maxSize:    15,
			data:       make([]byte, 16),
			bufferSize: 15,
			fileSize:   1,
		},
		{
			maxSize:    15,
			data:       make([]byte, 20),
			bufferSize: 15,
			fileSize:   5,
		},
		{
			maxSize:    20,
			data:       make([]byte, 1<<20),
			bufferSize: 20,
			fileSize:   1<<20 - 20,
		},
	}

	assert := assert.New(t)

	for i, tt := range tests {
		b := NewBuffer(tt.maxSize)

		n, err := b.Write(tt.data)
		if !assert.Nilf(err, "Test #%d: error during Write()", i+1) {
			continue
		}

		b.Finish()

		// Checks
		if !assert.Equalf(len(tt.data), n, "Test #%d: not all data written", i+1) {
			continue
		}

		if !assert.Equalf(tt.bufferSize, b.buff.Len(), "Test #%d: buffer contains wrong amount of bytes", i+1) {
			continue
		}

		if len(tt.data) <= tt.maxSize {
			// Must skip file checks

			assert.Equalf("", b.filename, "Test #%d: buffer created excess file", i+1)

			continue
		}

		f, err := os.Open(b.filename)
		if !assert.Nilf(err, "Test #%d: can't open file %s", i+1, b.filename) {
			continue
		}

		fileSize := func() int {
			info, err := f.Stat()
			if err != nil {
				return 0
			}

			return int(info.Size())
		}()

		f.Close()

		if !assert.Equalf(tt.fileSize, fileSize, "Test #%d: buffer contains wrong amount of bytes", i+1) {
			continue
		}

		b.reset()
	}
}

func TestBuffer_WriteAndRead(t *testing.T) {
	tests := []struct {
		maxSize   int
		sliceSize int
		//
		data [][]byte
		//
		res []byte
	}{
		{
			maxSize:   20,
			sliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   1,
			sliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   5,
			sliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   5,
			sliceSize: 20,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   5,
			sliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   5,
			sliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   0,
			sliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:   0,
			sliceSize: 5,
			data:      [][]byte{},
			res:       nil,
		},
	}

	assert := assert.New(t)

	for i, tt := range tests {
		b := NewBuffer(tt.maxSize)

		for _, d := range tt.data {
			n, err := b.Write(d)
			assert.Nilf(err, "Test #%d: error during Write()", i+1)
			assert.Equalf(len(d), n, "Test #%d: not all data written", i+1)
		}

		b.Finish()

		var res []byte

		data := make([]byte, tt.sliceSize)
		for {
			n, err := b.Read(data)
			data = data[:n]
			res = append(res, data...)
			data = data[:cap(data)]

			if err == io.EOF {
				break
			}

			assert.Nilf(err, "Test #%d: error during Read()", i+1)
		}

		assert.Equal(tt.res, res, "Test #%d: wrong content was read", i+1)

		b.reset()
	}
}

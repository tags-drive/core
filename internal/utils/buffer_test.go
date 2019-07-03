package utils

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"testing"
	"time"

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
		maxSize       int
		readSliceSize int
		//
		data [][]byte
		//
		res []byte
	}{
		{
			maxSize:       20,
			readSliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       1,
			readSliceSize: 256,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 20,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 10,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       5,
			readSliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       0,
			readSliceSize: 5,
			data: [][]byte{
				[]byte("123"),
				[]byte("456"),
				[]byte("789"),
			},
			res: []byte("123456789"),
		},
		{
			maxSize:       0,
			readSliceSize: 5,
			data:          [][]byte{},
			res:           nil,
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

		res, err := readByChunks(b, tt.readSliceSize)
		if !assert.Nilf(err, "Test #%d: error during Read()", i+1) {
			continue
		}

		assert.Equal(tt.res, res, "Test #%d: wrong content was read", i+1)

		b.reset()
	}
}

func TestBuffer_FuzzTest(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		t.Run("", func(t *testing.T) {
			assert := assert.New(t)

			var (
				sliceSize      = rand.Intn(1<<10) + 1
				bufferSize     = rand.Intn(sliceSize * 2) // can be zero
				writeChunkSize = rand.Intn(sliceSize) + 1
				readChunkSize  = rand.Intn(sliceSize) + 1
			)

			defer func() {
				// Log only when test is failed
				if t.Failed() {
					t.Logf("sliceSize: %d; bufferSize: %d; writeChunkSize: %d; readChunkSize: %d\n",
						sliceSize, bufferSize, writeChunkSize, readChunkSize)
				}
			}()

			slice := make([]byte, sliceSize)
			for i := range slice {
				slice[i] = byte(rand.Intn(128))
			}

			b := NewBuffer(bufferSize)

			// Write slice by chunks
			err := writeByChunks(b, slice, writeChunkSize)
			if !assert.Nil(err, "error during Write()") {
				t.FailNow()
			}

			b.Finish()

			res, err := readByChunks(b, readChunkSize)
			if !assert.Nil(err, "error during Read()") {
				t.FailNow()
			}

			if !assert.Equal(slice, res, "wrong content was read") {
				t.FailNow()
			}

			b.reset()
		})
	}
}

func BenchmarkBuffer(b *testing.B) {
	benchs := []struct {
		description    string
		dataSize       int
		maxBufferSize  int
		writeChunkSize int
		readChunkSize  int
	}{
		{
			description:    "Buffer size is greater than data",
			dataSize:       1 << 20, // 1MB
			maxBufferSize:  2 << 20, // 2MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
		{
			description:    "Buffer size is equal to data",
			dataSize:       1 << 20, // 1MB
			maxBufferSize:  1 << 20, // 1MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
		{
			description:    "Buffer size is less than data",
			dataSize:       20 << 20, // 20MB
			maxBufferSize:  1 << 20,  // 1MB
			writeChunkSize: 1024,
			readChunkSize:  2048,
		},
	}

	for _, bench := range benchs {
		b.Run(bench.description, func(b *testing.B) {
			slice := make([]byte, bench.dataSize)
			for i := range slice {
				slice[i] = byte(rand.Intn(128))
			}

			b.ResetTimer()

			b.Run("bytes.Buffer", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buff := bytes.NewBuffer(nil)

					err := writeByChunks(buff, slice, bench.writeChunkSize)
					if err != nil {
						b.Fatalf("error during Write(): %s", err)
					}

					_, err = readByChunks(buff, bench.readChunkSize)
					if err != nil {
						b.Fatalf("error during Read(): %s", err)
					}
				}
			})

			b.Run("utils.Buffer", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					buff := NewBuffer(bench.maxBufferSize)

					err := writeByChunks(buff, slice, bench.writeChunkSize)
					if err != nil {
						b.Fatalf("error during Write(): %s", err)
					}

					buff.Finish()

					_, err = readByChunks(buff, bench.readChunkSize)
					if err != nil {
						b.Fatalf("error during Read(): %s", err)
					}
				}
			})
		})
	}

}

func writeByChunks(w io.Writer, source []byte, chunk int) error {
	// Write slice by chunks
	for i := 0; i < len(source); i += chunk {
		bound := i + chunk
		if bound > len(source) {
			bound = len(source)
		}

		_, err := w.Write(source[i:bound])
		if err != nil {
			return err
		}
	}

	return nil
}

func readByChunks(r io.Reader, chunk int) ([]byte, error) {
	var res []byte

	data := make([]byte, chunk)
	for {
		n, err := r.Read(data)
		data = data[:n]
		res = append(res, data...)
		data = data[:cap(data)]

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}
	}

	return res, nil
}

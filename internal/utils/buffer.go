package utils

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
)

const (
	tempFileNameLength = 5
)

var (
	ErrBufferFinished      = errors.New("buffer is finished")
	ErrBufferIsNotFinished = errors.New("buffer is not finished")
)

// Buffer is a buffer which can store data on a disk. It isn't thread-safe
type Buffer struct {
	maxInMemorySize int

	// buff is used to store data in memory
	buff bytes.Buffer

	// file is used to store data on a disk
	file     *os.File
	useFile  bool
	filename string

	writingFinished bool
	readingFinished bool
}

func NewBuffer(maxInMemorySize int) *Buffer {
	return &Buffer{
		maxInMemorySize: maxInMemorySize,
	}
}

// Write writes data into bytes.Buffer while size of the Buffer is less than maxInMemorySize.
// When size of Buffer is equal to maxInMemorySize, Write creates a temporary file and writes remaining data into this one.
func (b *Buffer) Write(data []byte) (n int, err error) {
	if b.writingFinished {
		return 0, ErrBufferFinished
	}

	if !b.useFile {
		if b.buff.Len()+len(data) <= b.maxInMemorySize {
			// Just write data into the buffer
			return b.buff.Write(data)
		}

		// We have to use a file. But fill the buffer at first
		bound := b.maxInMemorySize - b.buff.Len()
		n, err = b.buff.Write(data[:bound])
		if err != nil {
			return n, err
		}

		// Trim written bytes
		data = data[bound:]

		b.useFile = true

		// Create a file in TempDir
		b.filename = os.TempDir() + "/" + GenerateRandomString(tempFileNameLength) + ".tmp"
		b.file, err = os.Create(b.filename)
		if err != nil {
			return n, errors.Wrap(err, "can't create a temp file")
		}
	}

	// Write data into the file
	n1, err := b.file.Write(data)
	n += n1
	return n, err
}

// Finish closes a temporary file if needed. All calls of Buffer.Write() will return ErrBufferFinished after the call of this method
func (b *Buffer) Finish() {
	if b.file != nil {
		b.file.Close()

		b.file = nil
	}

	b.writingFinished = true
}

// Read reads data from bytes.Buffer or from a file. Finish() method must be called before read from Buffer.
// A temp file is deleted when Read() encounter io.EOF
//
// The first time when Read encounter io.EOF it returns <nil>. It is needed to remove a temp file from the disk
func (b *Buffer) Read(data []byte) (n int, err error) {
	if !b.writingFinished {
		return 0, ErrBufferIsNotFinished
	}

	if b.readingFinished {
		if b.file != nil {
			b.file.Close()
			os.Remove(b.filename)

			b.file = nil
			b.filename = ""
		}

		return 0, io.EOF
	}

	if b.buff.Len() != 0 {
		// Use buffer
		n, err = b.readFromBuffer(data)
		if err != nil {
			return n, err
		}

		if n < len(data) {
			if !b.useFile {
				b.readingFinished = true
				return n, nil
			}

			// Can use the file to fill the slice

			var n1 int

			temp := make([]byte, len(data)-n)
			n1, err = b.readFromFile(temp)
			temp = temp[:n1]
			copy(data[n:], temp)
			n += n1
		}

		return n, err
	}

	if b.useFile {
		n, err = b.readFromFile(data)
		if err == io.EOF {
			b.readingFinished = true

			// Reset error
			err = nil
		}

		return n, err
	}

	// Reaching this code means that we buffer is empty and we don't use a file. So, reading is finished

	b.readingFinished = true
	return 0, nil
}

func (b *Buffer) readFromBuffer(data []byte) (n int, err error) {
	return b.buff.Read(data)
}

func (b *Buffer) readFromFile(data []byte) (n int, err error) {
	if b.file == nil {
		b.file, err = os.Open(b.filename)
		if err != nil {
			return 0, errors.Wrapf(err, "can't open a temp file '%s'", b.filename)
		}
	}

	return b.file.Read(data)
}

// Test only function. It resets buffer and remove file if needed
func (b *Buffer) reset() {
	b.buff.Reset()

	if b.file != nil {
		b.file.Close()
	}

	if b.filename != "" {
		os.Remove(b.filename)
	}
}

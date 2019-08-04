package bs_test

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/minio/sio"
	"github.com/stretchr/testify/assert"

	bs "github.com/tags-drive/core/internal/storage/files/binary_storage"
)

const (
	testFolder          = "./test"
	dataFolder          = "./test/data"
	resizedImagesFolder = "./test/data/resized"
)

func TestDiskStorage_SaveFile(t *testing.T) {
	t.Run("No encryption", func(t *testing.T) {
		defer clear()

		assert := assert.New(t)

		cnf := bs.DiskStorageConfig{
			DataFolder:          dataFolder,
			ResizedImagesFolder: resizedImagesFolder,
			Encrypt:             false,
		}
		storage, err := bs.NewDiskStorage(cnf)
		if !assert.Nil(err) {
			t.FailNow()
		}

		const buffSize = 4096
		const id = 0

		t.Run("Usual file", func(t *testing.T) {
			original := generateRandomData(buffSize)
			cp := make([]byte, buffSize)
			copy(cp, original)

			buff := bytes.NewBuffer(cp)
			storage.SaveFile(buff, id, buffSize, false)

			path := dataFolder + "/" + strconv.Itoa(id)
			assert.True(checkFile(path, original, false, nil), "files are not equal")
		})

		t.Run("Resized image", func(t *testing.T) {
			original := generateRandomData(buffSize)
			cp := make([]byte, buffSize)
			copy(cp, original)

			buff := bytes.NewBuffer(cp)
			storage.SaveFile(buff, id, buffSize, true)

			path := resizedImagesFolder + "/" + strconv.Itoa(id)
			assert.True(checkFile(path, original, false, nil), "files are not equal")
		})
	})

	t.Run("With encryption", func(t *testing.T) {
		defer clear()

		assert := assert.New(t)

		cnf := bs.DiskStorageConfig{
			DataFolder:          dataFolder,
			ResizedImagesFolder: resizedImagesFolder,
			Encrypt:             true,
			PassPhrase:          generatePassPhrase(),
		}
		storage, err := bs.NewDiskStorage(cnf)
		if !assert.Nil(err) {
			assert.FailNow("can't create a new DiskStorage")
		}

		const buffSize = 4096
		const id = 1

		t.Run("Usual file", func(t *testing.T) {
			original := generateRandomData(buffSize)
			cp := make([]byte, buffSize)
			copy(cp, original)

			buff := bytes.NewBuffer(cp)
			storage.SaveFile(buff, id, buffSize, false)

			path := dataFolder + "/" + strconv.Itoa(id)
			assert.True(checkFile(path, original, true, cnf.PassPhrase[:]), "files are not equal")
		})

		t.Run("Resized image", func(t *testing.T) {
			original := generateRandomData(buffSize)
			cp := make([]byte, buffSize)
			copy(cp, original)

			buff := bytes.NewBuffer(cp)
			storage.SaveFile(buff, id, buffSize, true)

			path := resizedImagesFolder + "/" + strconv.Itoa(id)
			assert.True(checkFile(path, original, true, cnf.PassPhrase[:]), "files are not equal")
		})
	})
}

func TestDiskStorage_GetFile(t *testing.T) {
	// We can use this function for both tests because only DiskStorage configs differ
	generateAndRunTests := func(assert *assert.Assertions, storage *bs.DiskStorage) {
		tests := []struct {
			id         int
			resized    bool
			data       []byte
			fileExists bool
		}{
			// Files exist
			{id: 0, resized: false, data: generateRandomData(512), fileExists: true},
			{id: 0, resized: true, data: generateRandomData(128), fileExists: true},
			// Files don't exist
			{id: -1, resized: false, data: nil, fileExists: false},
			{id: -1, resized: true, data: nil, fileExists: false},
		}

		// Create files
		for i, tt := range tests {
			if !tt.fileExists {
				continue
			}

			cp := make([]byte, len(tt.data))
			copy(cp, tt.data)

			buff := bytes.NewBuffer(cp)
			err := storage.SaveFile(buff, tt.id, int64(len(tt.data)), tt.resized)
			if !assert.Nilf(err, "Test #%d: can't create a file", i) {
				assert.FailNow("can't create a file. Fail now")
			}
		}

		for i, tt := range tests {
			buff := &bytes.Buffer{}

			err := storage.GetFile(buff, tt.id, tt.resized)
			if tt.fileExists {
				if !assert.Nilf(err, "Test #%d: can't get file", i) {
					continue
				}
			} else {
				if !assert.NotNilf(err, "Test #%d: get non-existed file", i) {
					continue
				}
			}

			// Check file
			assert.Truef(bytes.Equal(tt.data, buff.Bytes()), "Test #%d: get wrong content", i)
		}
	}

	t.Run("No encryption", func(t *testing.T) {
		defer clear()

		assert := assert.New(t)

		cnf := bs.DiskStorageConfig{
			DataFolder:          dataFolder,
			ResizedImagesFolder: resizedImagesFolder,
			Encrypt:             false,
		}
		storage, err := bs.NewDiskStorage(cnf)
		if !assert.Nil(err) {
			assert.FailNow("can't create a new DiskStorage")
		}

		generateAndRunTests(assert, storage)
	})

	t.Run("With encryption", func(t *testing.T) {
		defer clear()

		assert := assert.New(t)

		cnf := bs.DiskStorageConfig{
			DataFolder:          dataFolder,
			ResizedImagesFolder: resizedImagesFolder,
			Encrypt:             true,
			PassPhrase:          generatePassPhrase(),
		}
		storage, err := bs.NewDiskStorage(cnf)
		if !assert.Nil(err) {
			t.FailNow()
		}

		generateAndRunTests(assert, storage)
	})
}

func TestDiskStorage_DeleteFile(t *testing.T) {
	assert := assert.New(t)

	storage, err := bs.NewDiskStorage(bs.DiskStorageConfig{
		DataFolder:          dataFolder,
		ResizedImagesFolder: resizedImagesFolder,
		Encrypt:             false,
	})
	if !assert.Nil(err) {
		assert.FailNow("can't create a new DiskStorage")
	}

	defer clear()

	tests := []struct {
		id         int
		resized    bool
		data       []byte
		fileExists bool
	}{
		// Files exist
		{id: 0, resized: false, data: generateRandomData(512), fileExists: true},
		{id: 0, resized: true, data: generateRandomData(128), fileExists: true},
	}

	// Create files
	for i, tt := range tests {
		if !tt.fileExists {
			continue
		}

		cp := make([]byte, len(tt.data))
		copy(cp, tt.data)

		buff := bytes.NewBuffer(cp)
		err := storage.SaveFile(buff, tt.id, int64(len(tt.data)), tt.resized)
		if !assert.Nilf(err, "Test #%d: can't create a file", i) {
			assert.FailNow("can't create a file. Fail now")
		}
	}

	// Delete files
	for i, tt := range tests {
		err := storage.DeleteFile(tt.id, tt.resized)
		assert.Nilf(err, "Test #%d: can't delete file", i)
	}
}

// clear removes test folders
func clear() {
	os.RemoveAll(testFolder)
}

// Check file compares data of file with passed path and byte slice
func checkFile(path string, original []byte, encrypt bool, passPhrase []byte) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()

	if !encrypt {
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return false
		}
		return bytes.Equal(original, data)
	}

	buff := &bytes.Buffer{}
	_, err = sio.Decrypt(buff, f, sio.Config{Key: passPhrase})
	if err != nil {
		return false
	}

	return bytes.Equal(original, buff.Bytes())
}

func generateRandomData(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}

func generatePassPhrase() [32]byte {
	slice := generateRandomData(32)
	var array [32]byte
	for i := 0; i < 32; i++ {
		array[i] = slice[i]
	}
	return array
}
